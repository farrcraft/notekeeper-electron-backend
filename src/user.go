package main

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// User is a single user in an account
type User struct {
	ID            uuid.UUID      `json:"id"`          // ID is the unique identifier of the user
	Email         string         `json:"email"`       // Email is the email address of the user
	Active        bool           `json:"-"`           // Active indicates whether the user is active or not
	Account       *Account       `json:"-"`           // Account is the account that the user belongs to
	Created       time.Time      `json:"created"`     // Created is the time when the user was created
	Updated       time.Time      `json:"updated"`     // Updated is the time when the user was last created
	AccountKey    []byte         `json:"account_key"` // AccountKey is the encrypted version of the account-level encryption key
	PassphraseKey []byte         `json:"-"`           // PassphraseKey is the key derived from the passphrase
	Salt          []byte         `json:"salt"`        // Salt is the unique salt for generating the passphrase key
	Shelves       []*Shelf       `json:"-"`           // Shelves is the set of shelves that belong to the user
	Logger        *logrus.Logger `json:"-"`
	DB            *bolt.DB       `json:"-"`
}

// NewUser creates a new user object
func NewUser(db *bolt.DB, logger *logrus.Logger, email string) *User {
	now := time.Now()
	user := &User{
		ID:      uuid.NewV4(),
		Email:   email,
		Active:  true,
		Created: now,
		Updated: now,
		DB:      db,
		Logger:  logger,
	}
	return user
}
