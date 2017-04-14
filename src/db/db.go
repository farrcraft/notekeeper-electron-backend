package db

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// DB Types
const (
	DBTypeMaster = iota + 1
	DBTypeAccount
	DBTypeUser
	DBTypeCollection
	DBTypeShelf
	DBTypeUnknown
)

// DB is a database instance
type DB struct {
	ID       uuid.UUID
	Type     int
	DB       *bolt.DB
	Path     string
	Filename string
	Logger   *logrus.Logger
}

// New creates a new database
func New(dbtype int, logger *logrus.Logger) *DB {
	if dbtype >= DBTypeUnknown {
		logger.Debug("Unrecognized DB type")
		return nil
	}
	db := &DB{
		Type:   dbtype,
		Logger: logger,
	}
	return db
}

// Open a database
func (db *DB) Open() {

}

// Close a database
func (db *DB) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}
