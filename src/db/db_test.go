package db

import (
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
	uuid "github.com/satori/go.uuid"
)

func TestDB(t *testing.T) {
	logger, hook := test.NewNullLogger()

	factory := NewFactory("/tmp", logger)

	db := factory.DB(TypeUnknown, uuid.Nil)
	if db != nil {
		t.Error("Expected nil db")
	}

	db = factory.DB(TypeAccount, uuid.Nil)
	if db == nil {
		t.Error("Expected new db")
	}

	hook.Reset()
}
