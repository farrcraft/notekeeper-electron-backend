package handler

import (
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
)

// GetUserShelves gets a list of user's shelves
func GetUserShelves(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getShelves(server, message, "user")
	return response, err
}

// GetAccountShelves gets a list of account's shelves
func GetAccountShelves(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getShelves(server, message, "account")
	return response, err
}

// CreateUserShelf saves a new user shelf
func CreateUserShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createShelf(server, message, "user")
	return response, err
}

// CreateAccountShelf saves a new account shelf
func CreateAccountShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createShelf(server, message, "account")
	return response, err
}

// SaveUserShelf saves an existing user shelf
func SaveUserShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveShelf(server, message, "user")
	return response, err
}

// SaveAccountShelf saves an existing account shelf
func SaveAccountShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveShelf(server, message, "account")
	return response, err
}

// DeleteUserShelf deletes an existing user shelf
func DeleteUserShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteShelf(server, message, "user")
	return response, err
}

// DeleteAccountShelf deletes an existing account shelf
func DeleteAccountShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteShelf(server, message, "account")
	return response, err
}
