package user

import (
	"flag"
	"os"
	"testing"

	"../db"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var harness struct {
	logger   *logrus.Logger
	registry *db.Registry
	hook     *test.Hook
	user     *User
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
func TestUser(t *testing.T) {
	setup()

	err := harness.registry.OpenMaster("./")
	if err != nil {
		t.Error("Failed to open master db - ", err)
	}

	testSaveUser(t)

	// can't test loading user until we've saved the user index and done a user lookup on the index
	testSaveIndex(t)
	testLookupIndex(t)

	testLoadUser(t)

	teardown(t)
}

func testSaveUser(t *testing.T) {
	email := "bob@notekeeper.io"
	userPassphrase := "password"
	accountID, err := uuid.NewV4()
	if err != nil {
		t.Error("Expected to create an account id - ", err)
	}

	newUser, err := New(harness.registry, harness.logger, accountID, email)
	if err != nil {
		t.Error("Expected to create new user - ", err)
	}

	// create a new db file for the user
	userDBKey := db.Key{
		ID:   newUser.ID,
		Type: db.TypeUser,
	}
	userDBHandle, err := harness.registry.NewHandle(userDBKey)
	if err != nil {
		t.Error("Expected to create user db - ", err)
	}

	// create a new account db file for the user index
	accountDBKey := db.Key{
		ID:   newUser.AccountID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := harness.registry.NewHandle(accountDBKey)
	if err != nil {
		t.Error("Expected to create account db - ", err)
	}

	err = newUser.CreateKeys([]byte(userPassphrase))
	if err != nil {
		t.Error("Expected to create user keys - ", err)
	}

	userDBHandle.EncryptedKey = newUser.UserKey
	accountDBHandle.EncryptedKey = newUser.AccountKey

	err = newUser.Save()
	if err != nil {
		t.Error("Expected to save new user - ", err)
	}

	// Keep a copy of this user for additional testing
	harness.user = newUser
}

func testSaveIndex(t *testing.T) {
	index := NewIndex(harness.user.AccountID, harness.registry, harness.logger)
	err := index.Save(harness.user, harness.user.PassphraseKey)
	if err != nil {
		t.Error("Expected to save index - ", err)
	}
}

func testLoadUser(t *testing.T) {
	email := "bob@notekeeper.io"
	userPassphrase := "password"
	newUser, err := New(harness.registry, harness.logger, harness.user.AccountID, email)
	if err != nil {
		t.Error("Expected to create new user - ", err)
	}

	// Normally user ID & Salt come from the index lookup() method
	newUser.ID = harness.user.ID
	newUser.Salt = harness.user.Salt
	err = newUser.Load(userPassphrase)
	if err != nil {
		t.Error("Expected to load user - ", err)
	}
}

func testLookupIndex(t *testing.T) {
	email := "bob@notekeeper.io"
	newUser, err := New(harness.registry, harness.logger, harness.user.AccountID, email)

	index := NewIndex(harness.user.AccountID, harness.registry, harness.logger)
	err = index.Lookup(newUser)
	if err != nil {
		t.Error("Expected to lookup user in index - ", err)
	}
}
