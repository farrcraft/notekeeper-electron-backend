package handler

import (
	"../account"
	"../api"
	"../codes"
	messages "../proto"
	"../rpc"

	"github.com/golang/protobuf/proto"
)

// GetAccountState returns the accessible state of the account
func GetAccountState(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.AccountStateResponse{
		Header:   rpc.NewResponseHeader(),
		SignedIn: false,
		Locked:   true,
		Exists:   false,
	}

	if server.Account != nil {
		response.SignedIn = true
		response.Locked = server.Account.IsLocked()
		response.Exists = true
	} else {
		count := account.MapCount(server.DBRegistry)
		if count > 0 {
			response.Exists = true
		}
	}

	return response, nil
}

// CreateAccount is the RPC method to create a new account
func CreateAccount(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.UserIdResponse{
		Header: rpc.NewResponseHeader(),
		User:   &messages.UserId{},
	}

	request := messages.CreateAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create account request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create the account
	api := api.New(server.DBRegistry, server.Logger)
	newAccount, err := api.CreateAccount(request.Name, request.Email, request.Passphrase)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	// make this the active account
	if err == nil {
		server.Account = newAccount
		server.UserState = rpc.UserStateSignedIn
	}

	response.User.AccountId = newAccount.ID.String()
	response.User.UserId = newAccount.ActiveUser.ID.String()

	return response, nil
}

// UnlockAccount is the RPC method to unlock the current account
func UnlockAccount(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.UnlockAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling unlock account request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	api := api.New(server.DBRegistry, server.Logger)
	err = api.UnlockAccount(server.Account, request.Passphrase)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	} else {
		server.UserState = rpc.UserStateSignedIn
	}

	return response, nil
}

// SigninAccount is the RPC method to sign in to an existing account
func SigninAccount(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.UserIdResponse{
		Header: rpc.NewResponseHeader(),
		User:   &messages.UserId{},
	}

	request := messages.SigninAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling signin account request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	api := api.New(server.DBRegistry, server.Logger)
	newAccount, err := api.SigninAccount(request.Name, request.Email, request.Passphrase)
	if err == nil {
		server.Account = newAccount
		response.User.AccountId = newAccount.ID.String()
		response.User.UserId = newAccount.ActiveUser.ID.String()
		server.UserState = rpc.UserStateSignedIn
	} else {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}

// SignoutAccount is the RPC method to sign out from the active account
func SignoutAccount(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	api := api.New(server.DBRegistry, server.Logger)
	err := api.SignoutAccount(server.Account)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}
	server.Account = nil
	server.UserState = rpc.UserStateSignedOut
	return response, nil
}

// LockAccount is the RPC method to lock the active account
func LockAccount(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	api := api.New(server.DBRegistry, server.Logger)
	err := api.LockAccount(server.Account)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}
	server.UserState = rpc.UserStateLocked
	return response, nil
}
