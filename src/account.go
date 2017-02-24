package main

import (
	pb "./proto"
)

type Account struct {
	Email string
	Id    string // UUID
	DB    *bolt.DB
}
