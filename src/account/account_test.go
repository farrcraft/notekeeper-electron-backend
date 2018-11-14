package account

import (
	"os"
	"testing"

	"../db"
	"../user"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestAccount(t *testing.T) {
	// Setup
	logger, hook := test.NewNullLogger()
	factory := db.NewFactory("./", logger)

	account := New(factory, logger, "test_account")
	if account.Name != "test_account" {
		t.Error("Expected account name to be test_account")
	}

	masterDB, err := factory.DB(db.TypeMaster, uuid.Nil)
	if err != nil {
		t.Error("Failed to create master db - ", err)
	}

	count := MapCount(factory)
	if count != 0 {
		t.Error("Expected account map count to be 0")
	}

	err = account.OpenAccountDb()
	if err != nil {
		t.Error("Expected to open account db - ", err)
	}

	userEmail := "bob@notekeeper.io"
	userPassphrase := "supersecret"
	user := user.New(factory, logger, account.ID, userEmail)
	err = user.CreateKeys([]byte(userPassphrase))
	if err != nil {
		t.Error("Expected to create user keys - ", err)
	}

	/*
		// save user
		err = user.Save()
		if err != nil {
			t.Error("Expected to save user")
		}
	*/
	account.Users = append(account.Users, user.Profile)
	account.ActiveUser = user

	err = account.Save()
	if err != nil {
		t.Error("Expected to save account")
	}

	count = MapCount(factory)
	if count != 1 {
		t.Error("Expected account map count to be 1")
	}

	// Teardown
	masterDB.Close()
	err = os.Remove(masterDB.Filename)
	if err != nil {
		t.Error("Failed to cleanup master db - ", err)
	}

	accountDB := factory.Find(db.TypeAccount, account.ID)
	accountDB.Close()
	err = os.Remove(accountDB.Filename)
	if err != nil {
		t.Error("Failed to cleanup account db - ", err)
	}

	hook.Reset()
}
