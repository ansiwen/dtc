package zmq

import (
	"fmt"
	"github.com/niclabs/dtcnode/message"
	"github.com/niclabs/tcrsa"
	"github.com/pebbe/zmq4"
	"log"
	"sync"
	"time"
)

// The domain of the ZMQ connection. This value must be the same in the server, or it will not work.
const TchsmDomain = "tchsm"

// The protocol used for the ZMQ connection. TCP is the best for this usage cases.
const TchsmProtocol = "tcp"

// This error represents a timeout situation
var TimeoutError = fmt.Errorf("timeout")

// ZMQ Structure represents a connection to a set of Nodes via ZMQ
// Messaging Protocol.
type Server struct {
	// config properties
	host    string        // Host of service
	port    uint16        // Port where server ROUTER socket runs on
	privKey string        // Private Key of node
	pubKey  string        // Public Key of node
	timeout time.Duration // Length of timeout in seconds
	// nodes
	nodes []*Node // List of connected nodes of type ZMQ
	// zmq context
	ctx    *zmq4.Context // ZMQ Context
	socket *zmq4.Socket  // ROUTER socket which receives the responses from the nodes
	// message related structs
	channel         chan *message.Message       // The channel where all the responses from router are sent.
	pendingMessages map[string]*message.Message // A map with requests without response. To know what messages I'm expecting.
	mutex           sync.Mutex                  // A mutex to operate the pendingMessages map.
	currentMessage  message.Type                // A label which indicates the operation the connection is doing right now. It avoids inconsistent states (i.e. ask for a type of resource and then collect another one).
}

// New returns a new ZMQ connection based in the configuration provided.
func New(config *Config) (conn *Server, err error) {
	context, err := zmq4.NewContext()
	if err != nil {
		return
	}
	if config.Timeout == 0 {
		config.Timeout = 10
	}
	conn = &Server{
		host:            config.Host,
		port:            config.Port,
		privKey:         config.PrivateKey,
		pubKey:          config.PublicKey,
		timeout:         time.Duration(config.Timeout) * time.Second,
		ctx:             context,
		channel:         make(chan *message.Message),
		pendingMessages: make(map[string]*message.Message),
	}
	nodes := make([]*Node, len(config.Nodes))
	for i := 0; i < len(config.Nodes); i++ {
		nodes[i] = &Node{
			host:   config.Nodes[i].Host,
			port:   config.Nodes[i].Port,
			pubKey: config.Nodes[i].PublicKey,
			conn:   conn,
		}
	}
	conn.nodes = nodes
	return
}

func (conn *Server) Open() (err error) {
	if conn.nodes == nil {
		return fmt.Errorf("not initialized. Use 'New' to create a new struct")
	}
	err = zmq4.AuthStart()
	if err != nil {
		return
	}

	// Add IPs and public keys from clients
	zmq4.AuthAllow(TchsmDomain, conn.getIPs()...)
	zmq4.AuthCurveAdd(TchsmDomain, conn.getPubKeys()...)

	// Create in
	in, err := conn.ctx.NewSocket(zmq4.ROUTER)
	if err != nil {
		return
	}

	if err = in.SetIdentity(conn.pubKey); err != nil {
		return
	}

	// Add our private key
	err = in.ServerAuthCurve(TchsmDomain, conn.privKey)
	if err != nil {
		return
	}

	// Bind
	log.Printf("binding our socket in %s", conn.getConnString())
	err = in.Bind(conn.getConnString())
	if err != nil {
		return
	}
	conn.socket = in

	// Now we connect to the clients
	for _, client := range conn.nodes {
		if err = client.connect(); err != nil {
			return
		}
	}
	// Start message polling
	go func() {
		log.Printf("Message polling started")
		for {

			rawMsg, err := conn.socket.RecvMessageBytes(0)
			log.Printf("New message received!")
			if err != nil {
				log.Printf("Error with new message: %v", err)
				continue
			}
			msg, err := message.FromBytes(rawMsg)
			log.Printf("Message is from node %s", msg.NodeID)
			if err != nil {
				log.Printf("Cannot parse messages: %s\n", err)
				continue
			}
			log.Printf("Sending message to channel")
			conn.channel <- msg
			log.Printf("Message sent to channel!")

		}
	}()

	return
}

func (conn *Server) SendKeyShares(keyID string, keys tcrsa.KeyShareList, meta *tcrsa.KeyMeta) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if len(keys) != len(conn.nodes) {
		return fmt.Errorf("number of keys is not equal to number of nodes")
	}
	if conn.currentMessage != message.None {
		return fmt.Errorf("cannot send key shares in a currentMessage state different to None")
	}
	for i, node := range conn.nodes {
		log.Printf("Sending key share to node in %s:%d", node.host, node.port)
		msg, err := node.sendKeyShare(keyID, keys[i], meta)
		if err != nil {
			return fmt.Errorf("error with node %d: %s", i, err)
		}
		conn.pendingMessages[msg.ID] = msg
	}
	conn.currentMessage = message.SendKeyShare
	return nil
}

func (conn *Server) AckKeyShares() error {
	conn.mutex.Lock()
	defer func() {
		conn.pendingMessages = make(map[string]*message.Message)
		conn.currentMessage = message.None
		conn.mutex.Unlock()
	}()
	if conn.currentMessage != message.SendKeyShare {
		return fmt.Errorf("cannot ack key shares in a currentMessage state different to sendKeyShare")
	}
	acked := 0
	log.Printf("timeout will be %s", conn.timeout)
	timer := time.After(conn.timeout)
	for {
		select {
		case msg := <-conn.channel:
			log.Printf("message received from node %s\n", msg.NodeID)
			if pending, exists := conn.pendingMessages[msg.ID]; exists {
				if err := msg.Ok(pending, 0); err != nil {
					log.Printf("error with message from node %s: %v\n", msg.ID, message.ErrorToString[msg.Error])
				}
				delete(conn.pendingMessages, msg.ID)
				acked++
				if acked == len(conn.nodes) {
					return nil
				}
			} else {
				log.Printf("unexpected message: %+v\n", msg)
			}
		case <-timer:
			return TimeoutError
		}
	}
}

func (conn *Server) AskForSigShares(keyID string, hash []byte) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if conn.currentMessage != message.None {
		return fmt.Errorf("cannot ask for sig shares in a currentMessage state different to None")
	}
	for _, node := range conn.nodes {
		msg, err := node.askForSigShare(keyID, hash)
		if err != nil {
			return fmt.Errorf("error asking sigshare with node %s: %s", node.getID(), err)
		}
		conn.pendingMessages[msg.ID] = msg
	}
	conn.currentMessage = message.AskForSigShare
	return nil
}

func (conn *Server) GetSigShares() (tcrsa.SigShareList, error) {
	conn.mutex.Lock()
	defer func() {
		conn.pendingMessages = make(map[string]*message.Message)
		conn.currentMessage = message.None
		conn.mutex.Unlock()
	}()
	if conn.currentMessage != message.AskForSigShare {
		return nil, fmt.Errorf("cannot get sig shares in a currentMessage state different to askForSigShare")
	}
	sigShares := make(tcrsa.SigShareList, 0)
	timer := time.After(conn.timeout)
	minLen := 1
L:
	for {
		select {
		case msg := <-conn.channel:
			if pending, exists := conn.pendingMessages[msg.ID]; exists {
				if err := msg.Ok(pending, minLen); err != nil {
					log.Printf("error with message: %s\n", err)
					break
				}
				// Remove message from pending list
				delete(conn.pendingMessages, msg.ID)
				sigShare, err := message.DecodeSigShare(msg.Data[0])
				if err != nil {
					log.Printf("corrupt key: %v\n", msg)
					// Ask for it again?
					node := conn.getNodeByID(msg.NodeID)
					newMsg, err := node.askForSigShare(string(pending.Data[0]), pending.Data[1])
					if err != nil {
						log.Printf("error asking signature share to node %s: %s\n", node.getID(), err)
					}
					// save it in pending
					conn.pendingMessages[newMsg.ID] = newMsg
				} else {
					sigShares = append(sigShares, sigShare)
					if len(sigShares) == len(conn.nodes) {
						log.Printf("all signature shares retrieved.\n")
						break L
					}
				}
			} else {
				log.Printf("unexpected message: %v\n", msg)
			}
		case <-timer:
			log.Printf("timeout: %d out of %d sigs retrieved\n", len(sigShares), len(conn.nodes))
			break L
		}
	}
	return sigShares, nil
}

func (conn *Server) Close() error {
	err := conn.socket.Disconnect(conn.getConnString())
	if err != nil {
		return err
	}
	// Try to connect to peers
	for _, node := range conn.nodes {
		err = node.disconnect()
		if err != nil {
			return err
		}
	}
	zmq4.AuthStop()
	return nil
}

func (conn *Server) getPubKeys() []string {
	pubKeys := make([]string, len(conn.nodes))
	for i, node := range conn.nodes {
		pubKeys[i] = node.pubKey
	}
	return pubKeys
}

func (conn *Server) getIPs() []string {
	ips := make([]string, len(conn.nodes))
	for i, node := range conn.nodes {
		ips[i] = node.host
	}
	return ips
}

func (conn *Server) getConnString() string {
	return fmt.Sprintf("%s://%s:%d", TchsmProtocol, conn.host, conn.port)
}

func (conn *Server) getNodeByID(id string) *Node {
	for _, node := range conn.nodes {
		if node.getID() == id {
			return node
		}
	}
	return nil
}
