package main

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User is a single user in an account
type User struct {
	ID      uuid.UUID `json:"id"`      // ID is the unique identifier of the user
	Email   string    `json:"email"`   // Email is the email address of the user
	Active  bool      `json:"-"`       // Active indicates whether the user is active or not
	Account *Account  `json:"-"`       // Account is the account that the user belongs to
	Created time.Time `json:"created"` // Created is the time when the user was created
	Updated time.Time `json:"updated"` // Updated is the time when the user was last created
	Shelves []*Shelf  `json:"-"`       // Shelves is the set of shelves that belong to the user
}

// NewUser creates a new user object
func NewUser() *User {
	now := time.Now()
	user := &User{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return user
}
