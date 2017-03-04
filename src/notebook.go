package main

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
)

// Notebook contains a set of notes
type Notebook struct {
	ID           uuid.UUID `json:"id"`             // ID is the unique identifier for this notebook
	Account      *Account  `json:"-"`              // Account is the account that the notebook belongs to
	Title        *Title    `json:"title"`          // Title is the title of the notebook
	Default      bool      `json:"default"`        // Default indicates whether this is the default notebook
	EncryptedKey []byte    `json:"encryption_key"` // EncryptedKey is the encrypted version of the notebook's encryption key
	Notes        []*Note   `json:"-"`              // Notes is the set of notes that belong to this notebook
	Tags         []*Tag    `json:"tags"`           // Tags is the set of tags assigned to this notebook
	Created      time.Time `json:"created"`        // Created is the time when the notebook was created
	Updated      time.Time `json:"updated"`        // Updated is the time when the notebook was last updated
	Locked       bool      `json:"locked"`         // Locked indicates whether the notebook can be modified
}

// NewNotebook creates a new notebook object
func NewNotebook(account *Account) *Notebook {
	now := time.Now()
	notebook := &Notebook{
		ID:      uuid.NewV4(),
		Account: account,
		Created: now,
		Updated: now,
	}
	return notebook
}

// Save saves a notebook to the database
func (notebook *Notebook) Save() error {
	notebook.Account.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("notebooks"))
		if err != nil {
			notebook.Account.Logger.Error("Error creating notebook bucket - ", err)
			return err
		}

		// serialize notebook data
		data, err := json.Marshal(notebook)
		if err != nil {
			notebook.Account.Logger.Error("Error marshaling notebook - ", err)
			return err
		}

		// retrieve the encryption key
		_, userKey := ExtractSalt(notebook.Account.ActiveUser.PassphraseKey)
		decryptedKey, err := Decrypt(userKey, notebook.EncryptedKey)
		if err != nil {
			notebook.Account.Logger.Error("Error retrieving notebook key - ", err)
			return err
		}

		// decrypted encryption key is a byte slice, but to encrypt we need the key to be an array
		var key [KeySize]byte
		copy(key[:], decryptedKey[0:32])

		// encrypt the data
		encryptedData, err := Encrypt(&key, data)
		if err != nil {
			notebook.Account.Logger.Error("Error encrypting notebook data - ", err)
			return err
		}

		// finally, save it
		err = bucket.Put(notebook.ID.Bytes(), encryptedData)
		if err != nil {
			notebook.Account.Logger.Error("Error saving notebook - ", err)
			return err
		}
		return nil
	})

	return nil
}
