package db

import (
	"fmt"
	"path/filepath"

	"../codes"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
)

// Factory creates DB objects
type Factory struct {
	DataPath string
	Logger   *logrus.Logger
}

// IndexEntry is the stub of an index record used for loading the encrypted key
type IndexEntry struct {
	EncryptedKey []byte `json:"encryption_key"`
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

// CreateHandle returns a new DB handle.
// If the db file does not already exist, it will be created. If there is an
// existing db file, then it will be opened.
func (factory *Factory) CreateHandle(key Key) (*Handle, error) {
	if !IsValidType(key.Type) {
		code := codes.New(codes.ScopeDB, codes.ErrorInvalidType)
		return nil, code
	}

	handle := &Handle{
		Info: Info{
			ID:   key.ID,
			Type: key.Type,
		},
		Logger: factory.Logger,
	}

	var err error
	if key.Type == TypeMaster {
		handle.Info.Filename = filepath.Join(factory.DataPath, MasterDbFile)
	} else {
		if handle.Info.ID == uuid.Nil {
			handle.Info.ID, err = uuid.NewV4()
			if err != nil {
				return nil, err
			}
		}
		dbFile := fmt.Sprint(handle.Info.ID.String(), ".db")
		handle.Info.Filename = filepath.Join(factory.DataPath, dbFile)
	}

	err = handle.Open()
	if err != nil {
		return nil, err
	}

	return handle, nil
}
