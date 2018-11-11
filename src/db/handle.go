package db

import (
	"time"

	"../codes"

	"github.com/Sirupsen/logrus"
	"go.etcd.io/bbolt"
)

// Handle to an open database instance
type Handle struct {
	Info         Info
	DB           *bbolt.DB
	EncryptedKey []byte
	Logger       *logrus.Logger
}

// Open a database
func (handle *Handle) Open() error {
	var err error
	handle.DB, err = bbolt.Open(handle.Info.Filename, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		handle.Logger.Debug("Error opening DB type [", handle.Info.Type, "] file [", handle.Info.Filename, "] - ", err)
		var scope codes.Scope
		switch handle.Info.Type {
		case TypeMaster:
			scope = codes.ScopeGeneral
		case TypeAccount:
			scope = codes.ScopeAccount
		case TypeUser:
			scope = codes.ScopeUser
		case TypeCollection:
			scope = codes.ScopeCollection
		case TypeShelf:
			scope = codes.ScopeShelf
		}
		code := codes.New(scope, codes.ErrorDbOpen)
		return code
	}
	return nil
}

// Close a database
func (handle *Handle) Close() {
	if handle.DB != nil {
		handle.DB.Close()
		handle.DB = nil
	}
}
