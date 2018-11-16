package db

import (
	"../codes"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Registry of databases
type Registry struct {
	Factory *Factory
	Handles []*Handle
	Master  *Handle
	Logger  *logrus.Logger
}

// NewRegistry returns a new registry object
func NewRegistry(logger *logrus.Logger) *Registry {
	registry := &Registry{
		Logger:  logger,
		Factory: nil, // Not allocated until Registry::OpenMaster() is called
	}

	return registry
}

// OpenMaster opens the master DB
func (registry *Registry) OpenMaster(path string) error {
	registry.Factory = NewFactory(path, registry.Logger)

	// This is the master index db
	// There are additional databases where actual notebook data is stored
	key := Key{
		ID:   uuid.Nil,
		Type: TypeMaster,
	}
	handle, err := registry.NewHandle(key)
	if err != nil {
		return err
	}
	registry.Logger.Info("Opened master db file [", handle.Info.Filename, "]")

	return nil
}

// GetHandle to an already open database
func (registry *Registry) GetHandle(key Key) (*Handle, error) {
	// db handle already opened?
	for _, handle := range registry.Handles {
		if handle.Info.ID == key.ID && handle.Info.Type == key.Type {
			return handle, nil
		}
	}

	code := codes.New(codes.ScopeDB, codes.ErrorMissingDB)
	return nil, code
}

// NewHandle creates a new database handle.
// The database file will be opened and the handle registered, but the client
// will be responsible for assigning the encryption key to the handle.
func (registry *Registry) NewHandle(key Key) (*Handle, error) {
	if registry.Factory == nil {
		code := codes.New(codes.ScopeDB, codes.ErrorDbOpen)
		return nil, code
	}

	// not open yet, need to open it
	handle, err := registry.Factory.CreateHandle(key)
	if err != nil {
		return nil, err
	}

	registry.Handles = append(registry.Handles, handle)
	if handle.Info.Type == TypeMaster {
		registry.Master = handle
	}

	return handle, nil
}

// CloseAll opened DBs
func (registry *Registry) CloseAll() error {
	err := registry.CloseAccountDBs()
	if err != nil {
		return err
	}

	if registry.Master != nil {
		err := registry.Master.Close()
		if err != nil {
			return err
		}
		registry.Master = nil
	}

	registry.Handles = nil

	return nil
}

// CloseAccountDBs closes everything except the master DB
func (registry *Registry) CloseAccountDBs() error {
	for _, handle := range registry.Handles {
		if handle.Info.Type != TypeMaster {
			err := handle.Close()
			if err != nil {
				return err
			}
		}
	}
	registry.Handles = nil
	registry.Handles = append(registry.Handles, registry.Master)

	return nil
}
