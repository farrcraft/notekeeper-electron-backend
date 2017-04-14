package db

import (
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
)

func TestDB(t *testing.T) {
	logger, hook := test.NewNullLogger()

	db := New(DBTypeUnknown, logger)
	if db != nil {
		t.Error("Expected nil db")
	}

	db = New(DBTypeAccount, logger)
	if db == nil {
		t.Error("Expected new db")
	}

	hook.Reset()
}
