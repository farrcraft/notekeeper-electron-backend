package main

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
)

// Notebook contains a set of notes
type Notebook struct {
	ID      uuid.UUID `json:"id"`      // ID is the unique identifier for this notebook
	Account *Account  `json:"-"`       // Account is the account that the notebook belongs to
	Title   *Title    `json:"title"`   // Title is the title of the notebook
	Default bool      `json:"default"` // Default indicates whether this is the default notebook
	Notes   []*Note   `json:"-"`       // Notes is the set of notes that belong to this notebook
	Tags    []*Tag    `json:"tags"`    // Tags is the set of tags assigned to this notebook
	Created time.Time `json:"created"` // Created is the time when the notebook was created
	Updated time.Time `json:"updated"` // Updated is the time when the notebook was last updated
	Locked  bool      `json:"locked"`  // Locked indicates whether the notebook can be modified
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
	// [FIXME] - encrypt requested notebook name value
	//ciphertext, _ := cryptosecretbox.CryptoSecretBox([]byte(request.Name), nonce, key)

	// [FIXME] - need to use account-specific DB
	notebook.Account.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("notebooks"))
		if err != nil {
			notebook.Account.Logger.Error("Error creating notebook bucket - ", err)
			return err
		}
		data, err := json.Marshal(notebook)
		if err != nil {
			notebook.Account.Logger.Error("Error marshaling notebook - ", err)
			return err
		}
		// [FIXME] - use encrypted data value
		err = bucket.Put(notebook.ID.Bytes(), data)
		if err != nil {
			notebook.Account.Logger.Error("Error saving notebook - ", err)
			return err
		}
		return nil
	})

	return nil
}
