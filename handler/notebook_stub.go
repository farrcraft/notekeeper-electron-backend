package handler

import (
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
)

// CreateUserNotebook is the RPC method to create a new user notebook
func CreateUserNotebook(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := createNotebook(server, message, "user", context)
	return response, err
}

// CreateAccountNotebook is the RPC method to create a new account notebook
func CreateAccountNotebook(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := createNotebook(server, message, "account", context)
	return response, err
}

// GetUserNotebooks gets all of the user notebooks
func GetUserNotebooks(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := getNotebooks(server, message, "user", context)
	return response, err
}

// GetAccountNotebooks gets all of the user notebooks
func GetAccountNotebooks(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := getNotebooks(server, message, "user", context)
	return response, err
}

// SaveUserNotebook saves an existing account notebook
func SaveUserNotebook(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := saveNotebook(server, message, "user", context)
	return response, err
}

// SaveAccountNotebook saves an existing user notebook
func SaveAccountNotebook(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := saveNotebook(server, message, "account", context)
	return response, err
}

// DeleteUserNotebook deletes a user notebook
func DeleteUserNotebook(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := deleteNotebook(server, message, "user", context)
	return response, err
}

// DeleteAccountNotebook deletes a account notebook
func DeleteAccountNotebook(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response, err := deleteNotebook(server, message, "account", context)
	return response, err
}
