package handler

import (
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
)

// GetUserTags gets a set of user tags
func GetUserTags(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getTags(server, message, "user")
	return response, err
}

// GetAccountTags gets a set of account tags
func GetAccountTags(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := getTags(server, message, "account")
	return response, err
}

// CreateUserTag creates a new user tag
func CreateUserTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createTag(server, message, "user")
	return response, err
}

// CreateAccountTag creates a new account tag
func CreateAccountTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := createTag(server, message, "account")
	return response, err
}

// SaveUserTag saves an existing user tag
func SaveUserTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveTag(server, message, "user")
	return response, err
}

// SaveAccountTag saves an existing account tag
func SaveAccountTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := saveTag(server, message, "account")
	return response, err
}

// DeleteUserTag deletes a user tag
func DeleteUserTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteTag(server, message, "user")
	return response, err
}

// DeleteAccountTag deletes a account tag
func DeleteAccountTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response, err := deleteTag(server, message, "account")
	return response, err
}
