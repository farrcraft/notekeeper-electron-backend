package main

import ()

const (
	USER_SHELF    = "user"
	ACCOUNT_SHELF = "account"
)

type Shelf struct {
	Title     *Title
	Id        string
	Type      string
	Default   bool
	Trash     bool
	Account   *Account
	User      *User
	Notebooks []*Notebook
	Tags      []*Tag
	Created   string
	Updated   string
	Locked    bool
}

func NewShelf() *Shelf {

}
