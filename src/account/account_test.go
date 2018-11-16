package account

import (
	"os"
	"testing"

	"../db"
	"../user"

	"github.com/sirupsen/logrus/hooks/test"
)

func TestAccount(t *testing.T) {
	// Setup
	logger, hook := test.NewNullLogger()
	registry := db.NewRegistry(logger)

	accountName := "bob's account"
	newAccount, err := New(registry, logger, accountName)
	if err != nil {
		t.Error("Expected to create new account - ", err)
	}
	if newAccount.Name != accountName {
		t.Error("Expected account name to be bob's account")
	}

	err = registry.OpenMaster("./")
	if err != nil {
		t.Error("Failed to open master db - ", err)
	}

	accountIndex := NewIndex(registry, logger)

	count := accountIndex.Count()
	if count != 0 {
		t.Error("Expected account index count to be 0")
	}

	accountDBKey := db.Key{
		ID:   newAccount.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := registry.NewHandle(accountDBKey)
	if err != nil {
		t.Error("Expected to open account db - ", err)
	}

	userEmail := "bob@notekeeper.io"
	userPassphrase := "supersecret"
	newUser, err := user.New(registry, logger, newAccount.ID, userEmail)
	if err != nil {
		t.Error("Expected to create new user - ", err)
	}

	err = newUser.CreateKeys([]byte(userPassphrase))
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
	userDBKey := db.Key{
		ID:   newUser.ID,
		Type: db.TypeUser,
	}
	userDBHandle, err := registry.NewHandle(userDBKey)
	if err != nil {
		t.Error("Expected to create user db handle - ", err)
	}

	// [FIXME] - save user mapping in user index in account db

	// save user
	err = newUser.Save()
	if err != nil {
		t.Error("Expected to save user - ", err)
	}

	// [FIXME] - save user profile in account db
	// right now users are just stored as part of the account profile data
	newAccount.Users = append(newAccount.Users, newUser.Profile)
	newAccount.ActiveUser = newUser

	newAccount.EncryptedKey = newUser.AccountKey
	accountDBHandle.EncryptedKey = newUser.AccountKey
	userDBHandle.EncryptedKey = newUser.UserKey

	err = newAccount.Save()
	if err != nil {
		t.Error("Expected to save account - ", err)
	}

	err = accountIndex.Save(newAccount)
	if err != nil {
		t.Error("Expected to save account index")
	}

	count = accountIndex.Count()
	if count != 1 {
		t.Error("Expected account index count to be 1")
	}

	// Teardown
	err = registry.Master.DB.Close()
	if err != nil {
		t.Error("Failed to close master db - ", err)
	}
	err = os.Remove(registry.Master.Info.Filename)
	if err != nil {
		t.Error("Failed to cleanup master db - ", err)
	}

	for _, handle := range registry.Handles {
		if handle.Info.Type != db.TypeMaster {
			err = handle.Close()
			if err != nil {
				t.Error("Failed to close db - ", err)
			}
			err = os.Remove(handle.Info.Filename)
			if err != nil {
				t.Error("Failed to cleanup db - ", err)
			}
		}
	}

	hook.Reset()
}
