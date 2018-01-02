package db

import (
	"../codes"
	"github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
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
	handle, err := registry.GetHandle(key)
	registry.Logger.Info("Opened master db file [", handle.Info.Filename, "]")
	if err != nil {
		return err
	}

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
func (registry *Registry) CloseAll() {
	registry.CloseAccountDBs()
	if registry.Master != nil {
		registry.Master.Close()
		registry.Master = nil
	}
}

// CloseAccountDBs closes everything except the master DB
func (registry *Registry) CloseAccountDBs() {
	for _, handle := range registry.Handles {
		handle.Close()
	}
	registry.Handles = nil
}

/*
// Open a DB file & load its encrypted key
func (factory *Factory) Open(key Key, parentKey Key, ownerKey Key, passphraseKey []byte) (*DB, error) {
	var bucketName string
	var parentDB *DB

	if parentKey.Type != TypeAccount && parentKey.Type != TypeUser {
		code := codes.New(codes.ScopeDB, codes.ErrorInvalidType)
		return nil, code
	}

	if key.Type == TypeShelf {
		bucketName = "shelf_index"
		// parent should always be a findable account or user db
		parentDB = factory.Find(parentKey.Type, parentKey.ID)
	} else if key.Type == TypeCollection {
		bucketName = "collection_index"
		parentDB = factory.Find(TypeShelf, parentKey.ID)
		if parentDB == nil {
			var err error
			parentDB, err = factory.DB(TypeShelf, parentKey.ID)
			if err != nil {
				return nil, err
			}
			// had to open parent fresh which means we need to find the parent's encrypted key too
			// we need to load the shelf record from either the user or account db
			ownerDB := factory.Find(ownerKey.Type, ownerKey.ID)
			encryptedKey, err := factory.LoadEncryptedKey(parentKey.ID, passphraseKey, []byte("shelf_index"), ownerDB)
			if err != nil {
				return nil, err
			}
			parentDB.EncryptedKey = encryptedKey
		}
	} else {
		// unsupported type
		code := codes.New(codes.ScopeDB, codes.ErrorInvalidType)
		return nil, code
	}

	encryptedKey, err := factory.LoadEncryptedKey(key.ID, passphraseKey, []byte(bucketName), parentDB)
	if err != nil {
		return nil, err
	}

	db, err := factory.DB(key.Type, key.ID)
	if err != nil {
		return nil, err
	}
	db.EncryptedKey = encryptedKey
	return db, nil
}
*/
