package db

import (
	"time"

	"../codes"

	"github.com/sirupsen/logrus"
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
		handle.Logger.Warn("Error opening DB type [", handle.Info.Type, "] file [", handle.Info.Filename, "] - ", err)
		scope := dbTypeToErrorCode(handle.Info.Type)
		code := codes.New(scope, codes.ErrorDbOpen)
		return code
	}
	return nil
}

// dbTypeToErrorCode takes a DB type and returns a error code scope
func dbTypeToErrorCode(dbType Type) codes.Scope {
	var scope codes.Scope
	switch dbType {
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
	return scope
}

// Close a database
func (handle *Handle) Close() error {
	if handle.DB != nil {
		err := handle.DB.Close()
		if err != nil {
			handle.Logger.Warn("Error closing DB type [", handle.Info.Type, "] file [", handle.Info.Filename, "] - ", err)
			scope := dbTypeToErrorCode(handle.Info.Type)
			code := codes.New(scope, codes.ErrorDbClose)
			return code
		}
		handle.DB = nil
	}
	return nil
}
