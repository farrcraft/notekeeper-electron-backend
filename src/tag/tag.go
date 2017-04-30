package tag

import (
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../title"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Scope indicates the scope of the tag
type Scope int

const (
	// ScopeUser indicates that a tag belongs to a single user
	ScopeUser Scope = iota
	// ScopeAccount indicates that a tag belongs to a whole account
	ScopeAccount
)

// Tag is used for assigning labels to various object types
type Tag struct {
	ID        uuid.UUID      `json:"id"` // ID is the unique identifier of the tag
	AccountID uuid.UUID      `json:"-"`
	UserID    uuid.UUID      `json:"-"`
	Title     *title.Title   `json:"title"` // Title is the title of the tag
	Created   time.Time      `json:"created"`
	Updated   time.Time      `json:"updated"`
	Scope     Scope          `json:"scope"`
	DBFactory *db.Factory    `json:"-"` // DBFactory provides access to the database
	Logger    *logrus.Logger `json:"-"` // Logger is the logging facility
}

// New creates a new tag object
func New(title *title.Title, scope Scope, dbFactory *db.Factory, logger *logrus.Logger) *Tag {
	now := time.Now()
	tag := &Tag{
		ID:        uuid.NewV4(),
		Title:     title,
		Created:   now,
		Updated:   now,
		Scope:     scope,
		DBFactory: dbFactory,
		Logger:    logger,
	}
	return tag
}

func (tag *Tag) getDB() *db.DB {
	var dbType db.Type
	var id uuid.UUID
	if tag.Scope == ScopeUser {
		dbType = db.TypeUser
		id = tag.UserID
	} else {
		dbType = db.TypeAccount
		id = tag.AccountID
	}
	db := tag.DBFactory.Find(dbType, id)
	return db
}

// Save a tag to the DB
func (tag *Tag) Save(passphraseKey []byte) error {
	db := tag.getDB()
	err := db.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("tags"))
		if err != nil {
			tag.Logger.Debug("Error creating tag bucket - ", err)
			code := codes.New(codes.ScopeTag, codes.ErrorCreateBucket)
			return code
		}

		// serialize tag data
		data, err := json.Marshal(tag)
		if err != nil {
			tag.Logger.Debug("Error marshaling tag - ", err)
			code := codes.New(codes.ScopeTag, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		c := crypto.New(tag.Logger)
		decryptedKey, err := c.Open(passphraseKey, db.EncryptedKey)
		if err != nil {
			tag.Logger.Debug("Error retrieving tag key - ", err)
			code := codes.New(codes.ScopeTag, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := c.Seal(decryptedKey, data)
		if err != nil {
			tag.Logger.Debug("Error encrypting tag data - ", err)
			code := codes.New(codes.ScopeTag, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(tag.ID.Bytes(), encryptedData)
		if err != nil {
			tag.Logger.Debug("Error writing tag - ", err)
			code := codes.New(codes.ScopeTag, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		tag.Logger.Debug("Error saving tag - err")
		code := codes.New(codes.ScopeTag, codes.ErrorSave)
		return code
	}

	return nil
}

// LoadAll of the tags from an account or user DB
func (tag *Tag) LoadAll(passphraseKey []byte) ([]*Tag, error) {
	var tags []*Tag
	tagDB := tag.getDB()
	c := crypto.New(tag.Logger)
	tagKey, err := c.Open(passphraseKey, tagDB.EncryptedKey)
	if err != nil {
		tag.Logger.Debug("Error opening tag key - ", err)
		code := codes.New(codes.ScopeTag, codes.ErrorOpenKey)
		return tags, code
	}

	err = tagDB.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("tags"))
		if bucket == nil {
			tag.Logger.Debug("tag bucket does not exist")
			code := codes.New(codes.ScopeTag, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newTag := &Tag{
				DBFactory: tag.DBFactory,
				Logger:    tag.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(tagKey, value)
			if err != nil {
				tag.Logger.Debug("Error decrypting tag data - ", err)
				code := codes.New(codes.ScopeTag, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newTag)
			if err != nil {
				tag.Logger.Debug("Error decoding tag json - ", err)
				code := codes.New(codes.ScopeTag, codes.ErrorDecode)
				return code
			}

			tags = append(tags, newTag)
		}

		return nil
	})

	return tags, err
}

// Delete a tag
func (tag *Tag) Delete(passphraseKey []byte) error {
	tagDB := tag.getDB()
	err := tagDB.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tags"))
		if bucket == nil {
			tag.Logger.Debug("tag bucket does not exist")
			code := codes.New(codes.ScopeTag, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(tag.ID.Bytes())
		if err != nil {
			tag.Logger.Debug("Error deleting tag - ", err)
			code := codes.New(codes.ScopeTag, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		tag.Logger.Debug("Error deleting tag - ", err)
		code := codes.New(codes.ScopeTag, codes.ErrorDelete)
		return code
	}

	return nil
}
