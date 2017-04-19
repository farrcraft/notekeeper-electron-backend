package rpc

import (
	"../account"
	"../api"
	"../codes"
	messages "../proto"

	"github.com/golang/protobuf/proto"
)

// GetAccountState returns the accessible state of the account
func GetAccountState(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.AccountStateResponse{
		Header:   newResponseHeader(),
		SignedIn: false,
		Locked:   true,
		Exists:   false,
	}

	if rpc.Account != nil {
		response.SignedIn = true
		response.Locked = rpc.Account.IsLocked()
		response.Exists = true
	} else {
		count := account.MapCount(rpc.DBFactory)
		if count > 0 {
			response.Exists = true
		}
	}

	return response, nil
}

// CreateAccount is the RPC method to create a new account
func CreateAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.UserIdResponse{
		Header: newResponseHeader(),
		User:   &messages.UserId{},
	}

	request := messages.CreateAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling create account request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create the account
	api := api.New(rpc.DBFactory, rpc.Logger)
	newAccount, err := api.CreateAccount(request.Name, request.Email, request.Passphrase)
	if err != nil {
		setInternalError(response.Header, err)
		return response, nil
	}

	// make this the active account
	if err == nil {
		rpc.Account = newAccount
	}

	response.User.AccountId = newAccount.ID.String()
	response.User.UserId = newAccount.ActiveUser.ID.String()

	return response, nil
}

// UnlockAccount is the RPC method to unlock the current account
func UnlockAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	request := messages.UnlockAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling unlock account request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	api := api.New(rpc.DBFactory, rpc.Logger)
	err = api.UnlockAccount(rpc.Account, request.Passphrase)
	if err != nil {
		setInternalError(response.Header, err)
	}

	return response, nil
}

// SigninAccount is the RPC method to sign in to an existing account
func SigninAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.UserIdResponse{
		Header: newResponseHeader(),
		User:   &messages.UserId{},
	}

	request := messages.SigninAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling signin account request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	api := api.New(rpc.DBFactory, rpc.Logger)
	newAccount, err := api.SigninAccount(request.Name, request.Email, request.Passphrase)
	if err == nil {
		rpc.Account = newAccount
		response.User.AccountId = newAccount.ID.String()
		response.User.UserId = newAccount.ActiveUser.ID.String()
	}
	if err != nil {
		setInternalError(response.Header, err)
	}

	return response, nil
}

// SignoutAccount is the RPC method to sign out from the active account
func SignoutAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	api := api.New(rpc.DBFactory, rpc.Logger)
	err := api.SignoutAccount(rpc.Account)
	if err != nil {
		setInternalError(response.Header, err)
	}
	rpc.Account = nil

	return response, nil
}

// LockAccount is the RPC method to lock the active account
func LockAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	api := api.New(rpc.DBFactory, rpc.Logger)
	err := api.LockAccount(rpc.Account)
	if err != nil {
		setInternalError(response.Header, err)
	}

	return response, nil
}
