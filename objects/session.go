package objects

/*
#include "../criptoki/pkcs11go.h"
*/
import "C"
import (
	"bytes"
	"dtcmaster/network/zmq/message"
	"encoding/binary"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/niclabs/tcrsa"
	"reflect"
	"sync"
	"unsafe"
)

type CSessionHandle = C.CK_SESSION_HANDLE
type CSessionInfoPointer = C.CK_SESSION_INFO_PTR
type CState = C.CK_STATE
type CFlags = C.CK_FLAGS

const AttrTypeKeyHandler = 1 << 31
const AttrTypeKeyMeta = 1 << 31 + 1

type Session struct {
	sync.Mutex
	Slot            *Slot
	Handle          CSessionHandle
	flags           CFlags
	KeyMetaInfo     tcrsa.KeyMeta
	findInitialized bool
	refreshedToken  bool
	foundObjects    []CCryptoObjectHandle
	signInitialized bool
}

type Sessions map[CSessionHandle]*Session

var SessionHandle = CSessionHandle(0)

func NewSession(flags C.CK_FLAGS, currentSlot *Slot) *Session {
	SessionHandle++
	return &Session{
		Slot:   currentSlot,
		Handle: SessionHandle,
		flags:  flags,
	}
}

func (session *Session) GetHandle() CSessionHandle {
	return session.Handle
}

func (session *Session) GetCurrentSlot() *Slot {
	return session.Slot
}

func (session *Session) GetInfo(pInfo CSessionInfoPointer) error {
	if pInfo != nil {
		state, err := session.GetState()
		if err != nil {
			return err
		}
		info := (C.CK_SESSION_INFO)(unsafe.Pointer(pInfo))
		info.slotID = C.CK_SLOT_ID(session.Slot.ID)
		info.state = C.CK_STATE(state)
		info.flags = C.CK_FLAGS(session.flags)
		return nil

	} else {
		return NewError("Session.GetSessionInfo", "got NULL pointer", C.CKR_ARGUMENTS_BAD)
	}
}

// Saves an object and sets its handle.
func (session *Session) CreateObject(attrs Attributes) (*CryptoObject, error) {
	if attrs == nil {
		return nil, NewError("Session.CreateObject", "got NULL pointer", C.CKR_ARGUMENTS_BAD)
	}

	isTokenAttr, err := attrs.GetAttributeByType(C.CKA_TOKEN)
	if err != nil {
		return nil, NewError("Session.CreateObject", "is_token attr not defined", C.CKR_ARGUMENTS_BAD)
	}

	isToken := uint8(isTokenAttr.Value[0]) != 0
	var objType CryptoObjectType

	if isToken {
		objType = TokenObject
	} else {
		objType = SessionObject
	}

	object := &CryptoObject{
		Type: objType,
		Attributes: attrs,
	}

	token := session.Slot.token
	isPrivate := C.CK_TRUE
	oClass := C.CKO_VENDOR_DEFINED
	keyType := C.CKK_VENDOR_DEFINED


	privAttr, err := object.Attributes.GetAttributeByType(C.CKA_PRIVATE)
	if err == nil && len(privAttr.Value) > 0 {
		isPrivate = C.CK_BBOOL(privAttr.Value[0]) == C.CK_TRUE
	}

	classAttr, err := object.Attributes.GetAttributeByType(C.CKA_CLASS)
	if err == nil && len(classAttr.Value) > 0 {
		oClass = C.CK_OBJECT_CLASS(classAttr.Value[0])
	}

	keyTypeAttr, err := object.Attributes.GetAttributeByType(C.CKA_KEY_TYPE)
	if err == nil && len(classAttr.Value) > 0 {
		keyType = C.CK_KEY_TYPE(keyTypeAttr.Value[0])
	}

	if isToken == C.CK_TRUE && session.isReadOnly() {
		return nil, NewError("Session.CreateObject", "session is read only", C.CKR_SESSION_READ_ONLY)
	}
	state, err := session.GetState()
	if err != nil {
		return nil, err
	}
	if !GetUserAuthorization(state, isToken, isPrivate, true) {
		return nil, NewError("Session.CreateObject", "user not logged in", C.CKR_USER_NOT_LOGGED_IN)
	}

	switch oClass {
	case C.CKO_PUBLIC_KEY, C.CKO_PRIVATE_KEY:
		if keyType == C.CKK_RSA {
			token.AddObject(object)
			err := session.GetCurrentSlot().Application.Database.SaveToken(token)
			if err != nil {
				return nil, NewError("Session.CreateObject", err.Error(), C.CKR_DEVICE_ERROR)
			}
			return object, nil
		} else {
			return nil, NewError("Session.CreateObject", "key type not supported yet", C.CKR_ATTRIBUTE_VALUE_INVALID)
		}
	}
	return nil, NewError("Session.CreateObject", "object class not supported yet", C.CKR_ATTRIBUTE_VALUE_INVALID)
	// TODO: Verificar que los objetos sean válidos
}

func (session *Session) DestroyObject(hObject CCryptoObjectHandle) error {
	token, err := session.Slot.GetToken()
	if err != nil {
		return err
	}
	if object, err := token.GetObject(hObject); err != nil {
		return err
	} else {
		attr := object.FindAttribute(C.CKA_VENDOR_DEFINED)
		if attr != nil {
			privateAttr := object.FindAttribute(C.CKA_PRIVATE)
			if privateAttr != nil {
				isPrivate := C.CK_BBOOL(privateAttr.Value[0]) == C.CK_TRUE
				if isPrivate {
					// TODO: Delete key shares from DTC core
				}
			}
		}
		_ = token.DeleteObject(hObject)
		err := session.GetCurrentSlot().Application.Database.SaveToken(token)
		if err != nil {
			return NewError("Session.DestroyObject", err.Error(), C.CKR_DEVICE_ERROR)
		}
		return nil
	}
}

func (session *Session) FindObjectsInit(pTemplate CAttrPointer, ulCount C.CK_ULONG) error {
	if session.findInitialized {
		return NewError("Session.FindObjectsInit", "operation already initialized", C.CKR_OPERATION_ACTIVE)
	}
	if pTemplate == nil {
		return NewError("Session.FindObjectsInit", "got NULL pointer", C.CKR_ARGUMENTS_BAD)
	}

	token, err := session.GetCurrentSlot().GetToken()
	if err != nil {
		return err
	}

	if uint64(ulCount) == 0 {
		session.foundObjects = make([]CCryptoObjectHandle, len(token.Objects))
		i := 0
		for handle, _ := range token.Objects {
			session.foundObjects[i] = handle
			i++
		}
	} else {
		session.foundObjects = make([]CCryptoObjectHandle, 0)
		for handle, object := range token.Objects {
			if object.Match(pTemplate, ulCount) {
				session.foundObjects = append(session.foundObjects, handle)
			}
		}
	}

	// Si no se encontro el objecto, recargar la base de datos y buscar de
	// nuevo, puede que el objeto haya sido creado por otra instancia.
	if ulCount != 0 && len(session.foundObjects) == 0 && !session.refreshedToken {
		session.refreshedToken = true
		slot := session.GetCurrentSlot()
		token, err := slot.GetToken()
		if err != nil {
			return err
		}
		db := slot.Application.Database
		newToken, err := db.GetToken(token.Label)
		if err != nil {
			return NewError("Session.DestroyObject", err.Error(), C.CKR_DEVICE_ERROR)
		}
		token.CopyState(newToken)
		slot.InsertToken(newToken)
		return session.FindObjectsInit(pTemplate, ulCount)
	}

	// TODO: Verificar permisos de acceso
	session.findInitialized = true
	return nil
}

func (session *Session) FindObjects(maxObjectCount C.CK_ULONG) ([]CCryptoObjectHandle, error) {
	if !session.findInitialized {
		return nil, NewError("Session.FindObjects", "operation not initialized", C.CKR_OPERATION_NOT_INITIALIZED)
	}
	limit := len(session.foundObjects)
	if int(maxObjectCount) >= limit {
		limit = int(maxObjectCount)
	}
	resul := session.foundObjects[:limit]
	session.foundObjects = session.foundObjects[limit:]
	return resul, nil
}

func (session *Session) FindObjectsFinal() error {
	if !session.findInitialized {
		return NewError("Session.FindObjectsFinal", "operation not initialized", C.CKR_OPERATION_NOT_INITIALIZED)
	} else {
		session.findInitialized = false
		session.refreshedToken = false
	}
	return nil
}

func (session *Session) GetObject(handle CCryptoObjectHandle) (*CryptoObject, error) {
	token, err := session.Slot.GetToken()
	if err != nil {
		return nil, err
	}
	object, err := token.GetObject(handle)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (session *Session) GetState() (CState, error) {
	switch session.Slot.token.GetSecurityLevel() {
	case SecurityOfficer:
		return C.CKS_RW_SO_FUNCTIONS, nil
	case User:
		if session.isReadOnly() {
			return C.CKS_RO_USER_FUNCTIONS, nil
		} else {
			return C.CKS_RW_USER_FUNCTIONS, nil
		}
	case Public:
		if session.isReadOnly() {
			return C.CKS_RO_PUBLIC_SESSION, nil
		} else {
			return C.CKS_RW_PUBLIC_SESSION, nil
		}
	}
	return 0, NewError("Session.GetState", "invalid security level", C.CK_ARGUMENTS_BAD)
}

func (session *Session) isReadOnly() bool {
	return (session.flags & C.CKF_RW_SESSION) != C.CKF_RW_SESSION
}

func (session *Session) Login(userType C.CK_USER_TYPE, pPin C.CK_UTF8CHAR_PTR, ulPinLen C.CK_ULONG) error {
	token, err := session.Slot.GetToken()
	if err != nil {
		return err
	}
	return token.Login(userType, pPin, ulPinLen)
}

func (session *Session) Logout() error {
	token, err := session.Slot.GetToken()
	if err != nil {
		return err
	}
	token.Logout()
	return nil
}

func (session *Session) GenerateKeyPair(mechanism *Mechanism, pkAttrs, skAttrs Attributes) (pkObject, skObject *CryptoObject, err error) {
	// TODO: Verify access permissions (in my defense, the original implementation didn't do that too)
	if mechanism == nil || pkAttrs == nil || skAttrs == nil { // maybe this should be 0?
		err = NewError("Session.GenerateKeyPair", "got NULL pointer", C.CKR_ARGUMENTS_BAD)
		return
	}

	bitSizeAttr, err := pkAttrs.GetAttributeByType(C.CKA_MODULUS_BITS)
	if err != nil {
		err = NewError("Session.GenerateKeyPair", "got NULL pointer", C.CKR_TEMPLATE_INCOMPLETE)
		return
	}


	bitSize := binary.LittleEndian.Uint64(bitSizeAttr.Value)

	switch mechanism.Type {
	case C.CKM_RSA_PKCS_KEY_PAIR_GEN:
		keyID := uuid.New().String()
		// TODO: check if this UUID had been used before (?)
		var keyMeta *tcrsa.KeyMeta
		var pk, sk Attributes
		keyMeta, err = session.Slot.Application.DTC.CreateNewKey(keyID, int(bitSize), nil)
		if err != nil {
			return
		}
		pk, err = createPublicKey(keyID, pkAttrs, keyMeta)
		if err != nil {
			return
		}
		pkObject, err = session.CreateObject(pk)
		if err != nil {
			return
		}
		sk, err = createPrivateKey(keyID, skAttrs, keyMeta)
		if err != nil {
			return
		}
		skObject, err = session.CreateObject(sk)
		if err != nil {
			return
		}
	default:
		return nil, nil, NewError("Session.GenerateKeyPair", "mechanism invalid", C.CKR_MECHANISM_INVALID)
	}
	return pkObject, skObject, nil
}

func (session *Session) SignInit(mechanism *Mechanism, hKey CCryptoObjectHandle) error {
	if session.signInitialized {
		return NewError("Session.SignInit", "operation active", C.CKR_OPERATION_ACTIVE)
	}
	keyObject, err := session.GetObject(hKey)
	if err != nil {
		return err
	}
	keyNameAttr := keyObject.FindAttribute(AttrTypeKeyHandler)
	if keyNameAttr == nil {
		return NewError("Session.SignInit", "object handle does not contain any key", C.CKR_ARGUMENTS_BAD)
	}
	keyMetaAttr := keyObject.FindAttribute(AttrTypeKeyMeta)
	if keyMetaAttr == nil {
		return NewError("Session.SignInit", "object handle does not contain any key metainfo", C.CKR_ARGUMENTS_BAD)
	}

	keyNameStr := string(keyNameAttr.Value)
	keyMeta, err := message.DecodeKeyMeta(keyNameAttr.Value)
	if err != nil {
		return NewError("Session.SignInit", "key metainfo is corrupt", C.CKR_ARGUMENTS_BAD)
	}
	signMechanism := mechanism.Get

	return nil
}

func createPublicKey(keyID string, pkAttrs Attributes, keyMeta *tcrsa.KeyMeta) (Attributes, error) {

	// This fields are defined in SoftHSM implementation
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_CLASS, []byte{C.CKO_PUBLIC_KEY}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_KEY_TYPE, []byte{C.CKK_RSA}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_KEY_GEN_MECHANISM, []byte{C.CKM_RSA_PKCS_KEY_PAIR_GEN}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_LOCAL, []byte{C.CK_TRUE}})

	// This fields are our defaults

	pkAttrs.SetIfUndefined(&Attribute{C.CKA_LABEL, nil})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_ID, nil})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_SUBJECT, nil})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_PRIVATE, []byte{C.CK_FALSE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_MODIFIABLE, []byte{C.CK_TRUE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_TOKEN, []byte{C.CK_FALSE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_DERIVE, []byte{C.CK_FALSE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_ENCRYPT, []byte{C.CK_TRUE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_VERIFY, []byte{C.CK_TRUE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_VERIFY_RECOVER, []byte{C.CK_TRUE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_WRAP, []byte{C.CK_TRUE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_TRUSTED, []byte{C.CK_FALSE}})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_START_DATE, make([]byte,8)})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_END_DATE, make([]byte,8)})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_MODULUS_BITS, nil})

	// E and N from PK

	eBytes := make([]byte, reflect.TypeOf(keyMeta.PublicKey.E).Size())
	binary.LittleEndian.PutUint64(eBytes, uint64(keyMeta.PublicKey.E))

	pkAttrs.SetIfUndefined(&Attribute{C.CKA_MODULUS, keyMeta.PublicKey.N.Bytes()})
	pkAttrs.SetIfUndefined(&Attribute{C.CKA_PUBLIC_EXPONENT, eBytes})

	// Custom Fields

	encodedKeyMeta, err := encodeKeyMeta(keyMeta)
	if err != nil {
		return nil, err
	}

	pkAttrs.SetIfUndefined(&Attribute{CAttrType(AttrTypeKeyHandler), []byte(keyID)})
	pkAttrs.SetIfUndefined(&Attribute{CAttrType(AttrTypeKeyMeta), encodedKeyMeta})

	return pkAttrs, nil
}

func createPrivateKey(keyID string,skAttrs Attributes, keyMeta *tcrsa.KeyMeta) (Attributes, error) {

	// This fields are defined in SoftHSM implementation
	skAttrs.SetIfUndefined(&Attribute{C.CKA_CLASS, []byte{C.CKO_PRIVATE_KEY}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_KEY_TYPE, []byte{C.CKK_RSA}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_KEY_GEN_MECHANISM, []byte{C.CKM_RSA_PKCS_KEY_PAIR_GEN}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_LOCAL, []byte{C.CK_TRUE}})

	// This fields are our defaults

	skAttrs.SetIfUndefined(&Attribute{C.CKA_LABEL, nil})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_ID, nil})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_SUBJECT, nil})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_PRIVATE, []byte{C.CK_FALSE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_MODIFIABLE, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_TOKEN, []byte{C.CK_FALSE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_DERIVE, []byte{C.CK_FALSE}})

	skAttrs.SetIfUndefined(&Attribute{C.CKA_WRAP_WITH_TRUSTED, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_ALWAYS_AUTHENTICATE, []byte{C.CK_FALSE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_SENSITIVE, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_ALWAYS_SENSITIVE, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_DECRYPT, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_SIGN, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_DECRYPT, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_SIGN_RECOVER, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_UNWRAP, []byte{C.CK_TRUE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_EXTRACTABLE, []byte{C.CK_FALSE}})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_NEVER_EXTRACTABLE, []byte{C.CK_TRUE}})

	skAttrs.SetIfUndefined(&Attribute{C.CKA_START_DATE, make([]byte,8)})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_END_DATE, make([]byte,8)})

	// E and N from PK

	eBytes := make([]byte, reflect.TypeOf(keyMeta.PublicKey.E).Size())
	binary.LittleEndian.PutUint64(eBytes, uint64(keyMeta.PublicKey.E))

	skAttrs.SetIfUndefined(&Attribute{C.CKA_MODULUS, keyMeta.PublicKey.N.Bytes()})
	skAttrs.SetIfUndefined(&Attribute{C.CKA_PUBLIC_EXPONENT, eBytes})

	// Custom Fields

	encodedKeyMeta, err := encodeKeyMeta(keyMeta)
	if err != nil {
		return nil, err
	}

	skAttrs.SetIfUndefined(&Attribute{CAttrType(AttrTypeKeyHandler), []byte(keyID)})
	skAttrs.SetIfUndefined(&Attribute{CAttrType(AttrTypeKeyMeta), encodedKeyMeta})

	return skAttrs, nil
}

func encodeKeyMeta(meta *tcrsa.KeyMeta) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(meta); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}


func GetUserAuthorization(state CState, isToken, isPrivate C.CK_BBOOL, userAction bool) bool {
	switch state {
	case C.CKS_RW_SO_FUNCTIONS:
		return isPrivate == C.CK_FALSE
	case C.CKS_RW_USER_FUNCTIONS:
		return true
	case C.CKS_RO_USER_FUNCTIONS:
		if isToken == C.CK_TRUE {
			return !userAction
		} else {
			return true
		}
	case C.CKS_RW_PUBLIC_SESSION:
		return isPrivate == C.CK_FALSE
	case C.CKS_RO_PUBLIC_SESSION:
		return false
	}
	return false
}
