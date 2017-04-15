package db

import (
	"fmt"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

const (
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
)

// Factory creates DB objects
type Factory struct {
	DataPath string
	DBs      []*DB
	Logger   *logrus.Logger
}

// NewFactory creates a new database factory
func NewFactory(path string, logger *logrus.Logger) *Factory {
	cleanPath := filepath.Clean(path)

	factory := &Factory{
		DataPath: cleanPath,
		Logger:   logger,
	}
	return factory
}

// DB finds an existing runtime DB object or creates a new one
func (factory *Factory) DB(dbType int, id uuid.UUID) *DB {
	if dbType >= TypeUnknown {
		factory.Logger.Debug("Unrecognized DB type")
		return nil
	}

	db := factory.Find(dbType, id)
	if db != nil {
		return db
	}

	db = &DB{
		ID:     id,
		Type:   dbType,
		Logger: factory.Logger,
	}

	if dbType == TypeMaster {
		db.Filename = filepath.Join(factory.DataPath, MasterDbFile)
	} else {
		if id == uuid.Nil {
			db.ID = uuid.NewV4()
		}
		dbFile := fmt.Sprint(db.ID.String(), ".db")
		db.Filename = filepath.Join(factory.DataPath, dbFile)
	}

	factory.DBs = append(factory.DBs, db)

	return db
}

// Find an existing DB (runtime only, not on disk)
func (factory *Factory) Find(dbType int, id uuid.UUID) *DB {
	for _, db := range factory.DBs {
		if db.Type == dbType && db.ID == id {
			return db
		}
	}
	return nil
}

// CloseAll opened DBs
func (factory *Factory) CloseAll() {
	for _, db := range factory.DBs {
		db.Close()
	}
}

// CloseAccountDBs closes everything except the master DB
func (factory *Factory) CloseAccountDBs() {
	var master []*DB
	for _, db := range factory.DBs {
		if db.Type == TypeMaster {
			master = append(master, db)
		} else {
			db.Close()
		}
	}
	factory.DBs = master
}
