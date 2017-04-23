package db

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	uuid "github.com/satori/go.uuid"
)

func TestDB(t *testing.T) {
	logger, hook := test.NewNullLogger()

	factory := NewFactory("./", logger)

	db, err := factory.DB(42, uuid.Nil)
	if db != nil {
		t.Error("Expected nil db")
	}
	if err == nil {
		t.Error("Expected non-nil error")
	}

	db, err = factory.DB(TypeAccount, uuid.Nil)
	if db == nil {
		t.Error("Expected new db")
	}
	if err != nil {
		t.Error("Expected nil error, but got ", err)
	}

	db.Close()
	err = os.Remove(db.Filename)
	if err != nil {
		t.Error("Failed to cleanup test db - ", err)
	}

	hook.Reset()
}
