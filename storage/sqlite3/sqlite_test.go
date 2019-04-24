package sqlite3

import (
	"dtcmaster/objects"
	"dtcmaster/storage"
	"encoding/json"
	"fmt"
	"testing"
)

const TestFirstTokenLabel = "TCBHSM"
const TestOtherTokenLabel = "LABEL2"

func initDB() (storage.TokenStorage, error) {
	db, err := GetDatabase(":memory:")
	if err != nil {
		return nil, fmt.Errorf("couldn't get database: %v", err)
	}
	if err := db.InitStorage(); err != nil {
		return nil, fmt.Errorf("couldn't init storage in database: %v", err)
	}
	return db, nil
}

func TestDB_InitStorage(t *testing.T) {
	_, err := initDB()
	if err != nil {
		t.Errorf("init_storage: %v", err)
	}
}

func TestDB_GetMaxHandle(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Errorf("init_storage: %v", err)
	}
	maxHandle, err := db.GetMaxHandle()
	if err != nil {
		t.Errorf("couldn't get max handle: %v", err)
	}
	if maxHandle != 0 {
		t.Errorf("max handle is not zero on an empty database (it is %d)", maxHandle)
	}
}


func TestDB_GetToken(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Errorf("init_storage: %v", err)
	}
	token, err := db.GetToken(TestFirstTokenLabel)
	if err != nil {
		t.Errorf("get_token: %v", err)
		return
	}
	if token.Label != TestFirstTokenLabel {
		t.Errorf("wrong token retrieved, its label should have been %s, but it is %s", TestFirstTokenLabel, token.Label)
	}
}

func TestDB_SaveToken(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Errorf("init_storage: %v", err)
	}
	newToken := &objects.Token{
		Label:TestOtherTokenLabel,
		Pin: "1234",
		SoPin: "1234",
		Objects: make(objects.CryptoObjects),
	}
	actualHandle, err := db.GetMaxHandle()
	if err != nil {
		t.Errorf("max_handle: %v", err)
	}
	newHandle := actualHandle + 1
	newToken.Objects[newHandle] = &objects.CryptoObject{
		Handle: newHandle,
		Attributes: make(objects.Attributes),
	}
	newAttrType := int64(0)
	newToken.Objects[newHandle].Attributes[newAttrType] = &objects.Attribute{
		Type: newAttrType,
		Value: []byte("hello_world"),
	}
	err = db.SaveToken(newToken)
	if err != nil {
		t.Errorf("save_token: %v", err)
		return
	}
	gotToken, err := db.GetToken(TestOtherTokenLabel)
	if err != nil {
		t.Errorf("get_token: %v", err)
		return
	}
	if !newToken.Equals(gotToken) {
		newTokenJson, _ := json.Marshal(newToken)
		gotTokenJson, _ := json.Marshal(gotToken)
		t.Errorf(`wrong token retrieved
	expected: %s
	result:   %s`, newTokenJson, gotTokenJson)
	}
}

func TestDB_CloseStorage(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Errorf("init_storage: %v", err)
	}
	err = db.CloseStorage()
	if err != nil {
		t.Errorf("close_storage: %v", err)
	}
}