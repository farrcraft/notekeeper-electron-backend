package db

import (
	"encoding/json"
	"fmt"

	"../codes"
	"../crypto"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
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

// GetHandle to an open database
func (registry *Registry) GetHandle(key Key, passphraseKey []byte) (*Handle, error) {
	// db handle already opened?
	for _, handle := range registry.Handles {
		if handle.Info.ID == key.ID && handle.Info.Type == key.Type {
			return handle, nil
		}
	}

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

	// need to load the encrypted key for the db
	// may need to resolve a parent db first. hierarchy chain is:
	// master db -> account db -> user db -> shelf -> collection

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

// LoadEncryptedKey loads the encrypted key from an index bucket
func (registry *Registry) LoadEncryptedKey(dbKey Key, passphraseKey []byte, handle *Handle) ([]byte, error) {
	var encryptedKey []byte
	var bucketName []byte
	bucketName = []byte(fmt.Sprint(TypeToStr(dbKey.Type), "_index"))
	err := handle.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			registry.Logger.Debug(bucketName, " bucket does not exist")
			code := codes.New(codes.ScopeDB, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()
		key, value := cursor.Seek(dbKey.ID.Bytes())
		if key == nil {
			registry.Logger.Debug("Error loading record from index [", bucketName, "]")
			code := codes.New(codes.ScopeDB, codes.ErrorLoad)
			return code
		}

		c := crypto.New(registry.Logger)
		encryptionKey, err := c.Open(passphraseKey, handle.EncryptedKey)
		if err != nil {
			registry.Logger.Debug("Error opening key - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorOpenKey)
			return code
		}

		// decrypt value
		decryptedData, err := c.Open(encryptionKey, value)
		if err != nil {
			registry.Logger.Debug("Error decrypting data - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorDecrypt)
			return code
		}

		entry := &IndexEntry{}
		err = json.Unmarshal(decryptedData, entry)
		if err != nil {
			registry.Logger.Debug("Error decoding json - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorDecode)
			return code
		}

		encryptedKey = entry.EncryptedKey
		return nil
	})
	return encryptedKey, err
}
