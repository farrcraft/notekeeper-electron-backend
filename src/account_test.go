package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/boltdb/bolt"
)

func TestAccount(t *testing.T) {
	// Setup
	logger, hook := test.NewNullLogger()
	masterDbFileName := "master_test.db"
	masterDB, err := bolt.Open(masterDbFileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Error("Failed to create master db - ", err)
	}

	account := NewAccount(masterDB, logger, "test_account")
	if account.Name != "test_account" {
		t.Error("Expected account name to be test_account")
	}

	err = account.OpenAccountDb()
	if err != nil {
		t.Error("Expected to open account db - ", err)
	}

	userEmail := "bob@notekeeper.io"
	userPassphrase := "supersecret"
	user := NewUser(account.DB, logger, userEmail)
	err = user.CreateKeys([]byte(userPassphrase))
	if err != nil {
		t.Error("Expected to create user keys - ", err)
	}

	// save user
	err = user.Save()
	if err != nil {
		t.Error("Expected to save user")
	}
	account.Users = append(account.Users, user.Profile)
	account.ActiveUser = user

	err = account.Save()
	if err != nil {
		t.Error("Expected to save account")
	}

	// Teardown
	masterDB.Close()
	err = os.Remove(masterDbFileName)
	if err != nil {
		t.Error("Failed to cleanup master db - ", err)
	}

	account.DB.Close()
	accountDbFileName := fmt.Sprint(account.ID.String(), ".db")
	err = os.Remove(accountDbFileName)
	if err != nil {
		t.Error("Failed to cleanup account db - ", err)
	}

	hook.Reset()
}
