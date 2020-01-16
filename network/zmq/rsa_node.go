package zmq

import (
	"github.com/niclabs/dtcnode/v3/message"
	"github.com/niclabs/tcrsa"
)

func (node *Node) sendRSAKeyShare(id string, key *tcrsa.KeyShare, meta *tcrsa.KeyMeta) (*message.Message, error) {
	keyBinary, err := message.EncodeRSAKeyShare(key)
	if err != nil {
		return nil, err
	}
	metaBinary, err := message.EncodeRSAKeyMeta(meta)
	if err != nil {
		return nil, err
	}
	msg, err := message.NewMessage(message.SendRSAKeyShare, node.id(), []byte(id), keyBinary, metaBinary)
	if err != nil {
		return nil, err
	}
	_, err = node.socket.SendMessage(msg.GetBytesLists()...)
	if err != nil {
		return nil, err
	}
	return msg, nil

}

func (node *Node) getRSASigShare(id string, doc []byte) (msg *message.Message, err error) {
	msg, err = message.NewMessage(message.GetRSASigShare, node.id(), []byte(id), doc)
	if err != nil {
		return nil, err
	}
	if _, err := node.socket.SendMessage(msg.GetBytesLists()...); err != nil {
		return nil, err
	}
	return msg, nil
}

func (node *Node) deleteRSAKeyShare(id string) (*message.Message, error) {
	msg, err := message.NewMessage(message.DeleteRSAKeyShare, node.id(), []byte(id))
	if err != nil {
		return nil, err
	}
	_, err = node.socket.SendMessage(msg.GetBytesLists()...)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
