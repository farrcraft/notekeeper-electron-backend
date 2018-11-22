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
	testLoadUser(t)

	testSaveIndex(t)
	testLookupIndex(t)

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

	err = newUser.CreateKeys([]byte(userPassphrase))
	if err != nil {
		t.Error("Expected to create user keys - ", err)
	}

	userDBHandle.EncryptedKey = newUser.UserKey

	err = newUser.Save()
	if err != nil {
		t.Error("Expected to save new user - ", err)
	}
}

func testLoadUser(t *testing.T) {

}

func testSaveIndex(t *testing.T) {

}

func testLookupIndex(t *testing.T) {

}
