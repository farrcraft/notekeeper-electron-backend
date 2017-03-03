package main

import ()

type Notebook struct {
	Account *Account
	Id      string
	Title   *Title
	Default bool
	Notes   []*Note
	Tags    []*Tag
	Created string
	Updated string
	Locked  bool
}

func NewNotebook() *Notebook {
}
