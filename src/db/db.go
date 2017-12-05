package db

import (
	uuid "github.com/satori/go.uuid"
)

// Type indicates the type of DB
type Type int

// DB Types
const (
	TypeMaster Type = iota
	TypeAccount
	TypeUser
	TypeShelf
	TypeCollection
)

// Key to a DB
type Key struct {
	ID   uuid.UUID
	Type Type
}

// StrToType converts a string representation of a type to its native type value
func StrToType(typeName string) Type {
	var t Type
	switch typeName {
	case "master":
		t = TypeMaster
	case "account":
		t = TypeAccount
	case "user":
		t = TypeUser
	case "shelf":
		t = TypeShelf
	case "collection":
		t = TypeCollection
	}
	return t
}

// TypeToStr converts a type to its string representation
func TypeToStr(typeName Type) string {
	var name string
	switch typeName {
	case TypeMaster:
		name = "master"
	case TypeAccount:
		name = "account"
	case TypeUser:
		name = "user"
	case TypeShelf:
		name = "shelf"
	case TypeCollection:
		name = "collection"
	}
	return name
}

// IsValidType tests validity of a type value
func IsValidType(t Type) bool {
	if t != TypeMaster && t != TypeAccount && t != TypeUser && t != TypeShelf && t != TypeCollection {
		return false
	}
	return true
}
