package db

import (
	"time"

	"../codes"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// DB Types
const (
	TypeMaster = iota + 1
	TypeAccount
	TypeUser
	TypeCollection
	TypeShelf
	TypeUnknown
)

// DB is a database instance
type DB struct {
	ID       uuid.UUID
	Type     int
	DB       *bolt.DB
	Filename string
	Logger   *logrus.Logger
}

// Open a database
func (db *DB) Open() error {
	var err error
	db.DB, err = bolt.Open(db.Filename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		db.Logger.Debug("Error opening DB type [", db.Type, "] file [", db.Filename, "]")
		code := codes.New(codes.ErrorDbOpen)
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
