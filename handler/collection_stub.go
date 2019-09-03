package handler

import (
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
)

// GetUserCollections gets a list of a user's collections
func GetUserCollections(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getCollections(server, message, "user")
	return response, err
}

// GetAccountCollections gets a list of an account's collections
func GetAccountCollections(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getCollections(server, message, "account")
	return response, err
}

// CreateUserCollection saves a new user collection
func CreateUserCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createCollection(server, message, "user")
	return response, err
}

// CreateAccountCollection saves a new user collection
func CreateAccountCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createCollection(server, message, "account")
	return response, err
}

// SaveUserCollection updates a user collection
func SaveUserCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveCollection(server, message, "user")
	return response, err
}

// SaveAccountCollection updates a user collection
func SaveAccountCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveCollection(server, message, "account")
	return response, err
}

// DeleteUserCollection deletes a collection
func DeleteUserCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteCollection(server, message, "user")
	return response, err
}

// DeleteAccountCollection deletes a collection
func DeleteAccountCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteCollection(server, message, "account")
	return response, err
}
