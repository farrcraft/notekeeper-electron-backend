package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Account is the database holding one or more users and their collection of notes
type Account struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Users    []*User        `json:"users"`
	MasterDB *bolt.DB       `json:"-"`
	DB       *bolt.DB       `json:"-"`
	Shelves  []*Shelf       `json:"-"`
	Created  time.Time      `json:"created"`
	Updated  time.Time      `json:"updated"`
	Logger   *logrus.Logger `json:"-"`
}

// CreateAccountDb creates the database file for a new account
func (backend *Backend) CreateAccountDb(account *Account) error {
	fileName := fmt.Sprint(account.ID.String(), ".db")
	var err error
	account.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		backend.Logger.Error("Unable to open account DB [", fileName, "] - ", err)
		return errors.New("unable to open account DB")
	}
	backend.Account = account
	return nil
}

// NewAccount creates a new Account object
func NewAccount(db *bolt.DB, logger *logrus.Logger, name string) *Account {
	now := time.Now()
	account := &Account{
		ID:       uuid.NewV4(),
		Name:     name,
		MasterDB: db,
		Logger:   logger,
		Created:  now,
		Updated:  now,
	}
	return account
}

// Save saves the account to the database
func (account *Account) Save() error {
	// [FIXME] - missing an outer error handler
	// account data is stored in the master database
	account.MasterDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("accounts"))
		if err != nil {
			account.Logger.Error("Error creating accounts bucket - ", err)
			return err
		}
		data, err := json.Marshal(account)
		if err != nil {
			account.Logger.Error("Error marshaling account - ", err)
			return err
		}
		// [FIXME] - encrypt data value
		err = bucket.Put(account.ID.Bytes(), data)
		if err != nil {
			account.Logger.Error("Error saving account - ", err)
			return err
		}
		return nil
	})
	return nil
}
