package db

import (
	"time"

	"../codes"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Type indicates the type of DB
type Type int

// DB Types
const (
	TypeMaster Type = iota
	TypeAccount
	TypeUser
	TypeShelf
	TypeCollection
)

// DB is a database instance
type DB struct {
	ID           uuid.UUID
	Type         Type
	DB           *bolt.DB
	EncryptedKey []byte
	Filename     string
	Logger       *logrus.Logger
}

// Key to a DB
type Key struct {
	ID   uuid.UUID
	Type Type
}

// StrToType converts a string representation of a type to its native type value
func StrToType(typeName string) Type {
	var t Type
	switch typeName {
	case "master":
		t = TypeMaster
	case "account":
		t = TypeAccount
	case "user":
		t = TypeUser
	case "shelf":
		t = TypeShelf
	case "collection":
		t = TypeCollection
	}
	return t
}

// IsValidType tests validity of a type value
func IsValidType(t Type) bool {
	if t != TypeMaster && t != TypeAccount && t != TypeUser && t != TypeShelf && t != TypeCollection {
		return false
	}
	return true
}

// Open a database
func (db *DB) Open() error {
	var err error
	db.DB, err = bolt.Open(db.Filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		db.Logger.Debug("Error opening DB type [", db.Type, "] file [", db.Filename, "] - ", err)
		var scope codes.Scope
		switch db.Type {
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
func (db *DB) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
