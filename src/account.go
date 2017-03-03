package main

// Account is the database holding one or more users and their collection of notes
type Account struct {
	ID      string // UUID
	Name    string
	Users   []*User
	DB      *bolt.DB
	Shelves []*Shelf
	Created string
	Updated string
}

// NewAccount creates a new Account object
func NewAccount() *Account {
	account := &Account{
		ID: uuid.NewV4(),
	}
}
