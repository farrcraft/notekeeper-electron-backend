package handler

import (
	"notekeeper-electron-backend/account"
	"notekeeper-electron-backend/api"
	"notekeeper-electron-backend/codes"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
)

// GetAccountState returns the accessible state of the account
func GetAccountState(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
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
		accountIndex := account.NewIndex(server.DBRegistry, server.Logger)
		count := accountIndex.Count()
		if count > 0 {
			response.Exists = true
		}
		server.Logger.Debug("Account state counted [", count, "] accounts")
	}

	return response, nil
}

// CreateAccount is the RPC method to create a new account
func CreateAccount(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.UserIdResponse{
		Header: rpc.NewResponseHeader(),
		User:   &messages.UserId{},
	}

	request := messages.CreateAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling create account request - ", err)
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
	server.Account = newAccount
	server.UserState = rpc.UserStateSignedIn

	response.User.AccountId = newAccount.ID.String()
	response.User.UserId = newAccount.ActiveUser.ID.String()

	return response, nil
}

// SigninAccount is the RPC method to sign in to an existing account
func SigninAccount(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.UserIdResponse{
		Header: rpc.NewResponseHeader(),
		User:   &messages.UserId{},
	}

	request := messages.SigninAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling signin account request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	api := api.New(server.DBRegistry, server.Logger)
	newAccount, err := api.SigninAccount(request.Name, request.Email, request.Passphrase)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	server.Account = newAccount
	server.UserState = rpc.UserStateSignedIn

	response.User.AccountId = newAccount.ID.String()
	response.User.UserId = newAccount.ActiveUser.ID.String()

	return response, nil
}

// SignoutAccount is the RPC method to sign out from the active account
func SignoutAccount(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	api := api.New(server.DBRegistry, server.Logger)
	err := api.SignoutAccount(server.Account)
	server.Account = nil
	server.UserState = rpc.UserStateSignedOut
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}
	return response, nil
}

// LockAccount is the RPC method to lock the active account
func LockAccount(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	api := api.New(server.DBRegistry, server.Logger)
	err := api.LockAccount(server.Account)
	server.UserState = rpc.UserStateLocked
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}
	return response, nil
}

// UnlockAccount is the RPC method to unlock the current account
func UnlockAccount(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.UnlockAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling unlock account request - ", err)
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
