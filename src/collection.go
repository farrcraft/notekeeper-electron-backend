package main

import ()

type Collection struct {
	Id        string
	Title     *Title
	Notebooks []*Notebook
	Account   *Account
	User      *User
	Shelf     *Shelf
	Tags      []*Tag
	Created   string
	Updated   string
	Locked    bool
}
