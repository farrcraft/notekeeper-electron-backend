package account

import (
	"flag"
	"os"
	"testing"

	"../db"
	"../user"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var harness struct {
	logger   *logrus.Logger
	registry *db.Registry
	hook     *test.Hook
}

func TestMain(m *testing.M) {
	flag.Parse()

	exitCode := m.Run()

	os.Exit(exitCode)
}

func setup() {
	harness.logger, harness.hook = test.NewNullLogger()

	// null logger usually discards all log output, but for debugging we want to send all log output to stdout
	harness.logger.SetOutput(os.Stdout)
	harness.logger.SetLevel(logrus.DebugLevel)

	harness.registry = db.NewRegistry(harness.logger)
}

func teardown(t *testing.T) {
	// Teardown
	err := harness.registry.Master.DB.Close()
	if err != nil {
		t.Error("Failed to close master db - ", err)
	}
	err = os.Remove(harness.registry.Master.Info.Filename)
	if err != nil {
		t.Error("Failed to cleanup master db - ", err)
	}

	for _, handle := range harness.registry.Handles {
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

	harness.hook.Reset()
}

func TestAccount(t *testing.T) {
	setup()

	testCreateAccount(t)
	testSigninAccount(t)

	teardown(t)
}

func testCreateAccount(t *testing.T) {
	accountName := "bob's account"
	newAccount, err := New(harness.registry, harness.logger, accountName)
	if err != nil {
		t.Error("Expected to create new account - ", err)
	}
	if newAccount.Name != accountName {
		t.Error("Expected account name to be bob's account")
	}

	err = harness.registry.OpenMaster("./")
	if err != nil {
		t.Error("Failed to open master db - ", err)
	}

	accountIndex := NewIndex(harness.registry, harness.logger)

	count := accountIndex.Count()
	if count != 0 {
		t.Error("Expected account index count to be 0")
	}

	accountDBKey := db.Key{
		ID:   newAccount.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := harness.registry.NewHandle(accountDBKey)
	if err != nil {
		t.Error("Expected to open account db - ", err)
	}

	userEmail := "bob@notekeeper.io"
	userPassphrase := "supersecret"
	newUser, err := user.New(harness.registry, harness.logger, newAccount.ID, userEmail)
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
	userDBHandle, err := harness.registry.NewHandle(userDBKey)
	if err != nil {
		t.Error("Expected to create user db handle - ", err)
	}

	// save user mapping in user index in account db
	userIndex := user.NewIndex(newAccount.ID, harness.registry, harness.logger)
	userIndex.Save(newUser, newUser.PassphraseKey)

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
}

func testSigninAccount(t *testing.T) {
	// attempt to find the account (lookup)
	accountName := "bob's account"
	newAccount, err := New(harness.registry, harness.logger, accountName)
	if err != nil {
		t.Error("Expected to make new account object - ", err)
	}

	accountIndex := NewIndex(harness.registry, harness.logger)
	err = accountIndex.Lookup(newAccount)
	if err != nil {
		t.Error("Expected to lookup account in index - ", err)
	}

	// authenticate the user
	userEmail := "bob@notekeeper.io"
	userPassphrase := "supersecret"
	newUser, err := user.New(harness.registry, harness.logger, newAccount.ID, userEmail)
	if err != nil {
		t.Error("Expected to create new user - ", err)
	}

	// resolve the user id from the user map in the account db
	userIndex := user.NewIndex(newAccount.ID, harness.registry, harness.logger)
	err = userIndex.Lookup(newUser)
	if err != nil {
		t.Error("Expected to lookup user in index - ", err)
	}

	err = newUser.Load(userPassphrase)
	if err != nil {
		t.Error("Expected to load user - ", err)
	}

	// connect the user to the account & make it the active user
	newAccount.ActiveUser = newUser

	// set remaining encryption keys
	newAccount.EncryptedKey = newUser.AccountKey
	//accountDBHandle.EncryptedKey = newUser.AccountKey

	// load the account
	err = newAccount.Load()
	if err != nil {
		t.Error("Expected to load account - ", err)
	}

	if newAccount.Name != accountName {
		t.Error("Expected account name to match")
	}
}
