package db

import (
	uuid "github.com/satori/go.uuid"
)

// Info about a database
type Info struct {
	ID       uuid.UUID
	Type     Type
	Filename string
}
