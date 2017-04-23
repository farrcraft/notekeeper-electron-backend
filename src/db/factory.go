package db

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"../codes"
	"../crypto"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
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

// DB finds an existing runtime DB object or creates a new one
func (factory *Factory) DB(dbType Type, id uuid.UUID) *DB {
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
func (factory *Factory) Find(dbType Type, id uuid.UUID) *DB {
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

// Open a DB file & load its encrypted key
func (factory *Factory) Open(key Key, parentKey Key, ownerID uuid.UUID, passphraseKey []byte) (*DB, error) {
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
			parentDB = factory.DB(TypeShelf, parentKey.ID)
			// had to open parent fresh which means we need to find the parent's encrypted key too
			// we need to load the shelf record from either the user or account db
			ownerDB := factory.Find(parentKey.Type, ownerID)
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

	db := factory.DB(key.Type, key.ID)
	db.EncryptedKey = encryptedKey
	return db, nil
}

// LoadEncryptedKey loads the encrypted key from an index bucket
func (factory *Factory) LoadEncryptedKey(id uuid.UUID, passphraseKey []byte, bucketName []byte, db *DB) ([]byte, error) {
	var encryptedKey []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			factory.Logger.Debug(bucketName, " bucket does not exist")
			code := codes.New(codes.ScopeDB, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()
		key, value := cursor.Seek(id.Bytes())
		if key == nil {
			factory.Logger.Debug("Error loading record from index [", bucketName, "]")
			code := codes.New(codes.ScopeDB, codes.ErrorLoad)
			return code
		}

		encryptionKey, err := crypto.Open(passphraseKey, db.EncryptedKey)
		if err != nil {
			factory.Logger.Debug("Error opening key - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorOpenKey)
			return code
		}

		// decrypt value
		decryptedData, err := crypto.Open(encryptionKey, value)
		if err != nil {
			factory.Logger.Debug("Error decrypting data - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorDecrypt)
			return code
		}

		entry := &IndexEntry{}
		err = json.Unmarshal(decryptedData, entry)
		if err != nil {
			factory.Logger.Debug("Error decoding json - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorDecode)
			return code
		}

		encryptedKey = entry.EncryptedKey
		return nil
	})
	return encryptedKey, err
}
