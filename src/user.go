package main

import ()

type User struct {
	Id      string   `json:"id"`
	Email   string   `json:"email"`
	Active  bool     `json:"-"`
	Account *Account `json:"-"`
	Created string
	Updated string
	Shelves []*Shelf `json:"-"`
}

func NewUser() *User {

}
