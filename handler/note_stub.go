package handler

import (
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
)

// GetUserNotes is the RPC method to get a list of user notes
func GetUserNotes(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getNotes(server, message, "user")
	return response, err
}

// GetAccountNotes is the RPC method to get a list of account notes
func GetAccountNotes(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getNotes(server, message, "account")
	return response, err
}

// LoadUserNote is the RPC method to get a user note
func LoadUserNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := loadNote(server, message, "user")
	return response, err
}

// LoadAccountNote is the RPC method to get a account note
func LoadAccountNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := loadNote(server, message, "account")
	return response, err
}

// CreateUserNote is the RPC method to create a new user note
func CreateUserNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createNote(server, message, "user")
	return response, err
}

// CreateAccountNote is the RPC method to create a new account note
func CreateAccountNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createNote(server, message, "account")
	return response, err
}

// SaveUserNote is the RPC method to save an existing user note
func SaveUserNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveNote(server, message, "user")
	return response, err
}

// SaveAccountNote is the RPC method to save an existing account note
func SaveAccountNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveNote(server, message, "account")
	return response, err
}

// DeleteUserNote is the RPC method to delete a user note
func DeleteUserNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteNote(server, message, "user")
	return response, err
}

// DeleteAccountNote is the RPC method to delete a account note
func DeleteAccountNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteNote(server, message, "account")
	return response, err
}
