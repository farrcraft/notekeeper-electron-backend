package rpc

import (
	"../account"
	"../api"
	"../codes"
	"github.com/mitchellh/mapstructure"
)

type responseAccountState struct {
	SignedIn bool `json:"signed_in"`
	Locked   bool `json:"locked"`
	Exists   bool `json:"exists"`
}

// GetAccountState returns the accessible state of the account
func GetAccountState(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}
	payload := &responseAccountState{
		SignedIn: false,
		Locked:   true,
		Exists:   false,
	}
	if rpc.Account != nil {
		payload.SignedIn = true
		payload.Locked = rpc.Account.IsLocked()
		payload.Exists = true
	} else {
		count := account.MapCount(rpc.DB)
		if count > 0 {
			payload.Exists = true
		}
	}
	response.Payload = payload
	return response, nil
}

type requestCreateAccount struct {
	Name       string `mapstructure:"name"`
	Email      string `mapstructure:"email"`
	Passphrase string `mapstructure:"passphrase"`
}

type responseCreateAccount struct {
	ID string `json:"id"`
}

// CreateAccount is the RPC method to create a new account
func CreateAccount(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	var request requestCreateAccount
	err := mapstructure.Decode(message.Payload, &request)
	if err != nil {
		rpc.Logger.Debug("Error decoding create account request payload - ", err)
		response.Code = int(codes.ErrorCreateAccountDecode)
		response.Status = codes.StatusError
		return response, nil
	}

	// create the account
	newAccount, err := api.CreateAccount(rpc.DB, rpc.Logger, rpc.DataPath, request.Name, request.Email, request.Passphrase)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
		return response, nil
	}

	// make this the active account
	if err == nil {
		rpc.Account = newAccount
	}

	response.Payload = &responseCreateAccount{
		ID: newAccount.ID.String(),
	}

	return response, nil
}

type requestUnlockAccount struct {
	Passphrase string `mapstructure:"passphrase"`
}

// UnlockAccount is the RPC method to unlock the current account
func UnlockAccount(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	var request requestUnlockAccount
	err := mapstructure.Decode(message.Payload, &request)
	if err != nil {
		rpc.Logger.Debug("Error decoding unlock account request payload - ", err)
		response.Code = int(codes.ErrorUnlockAccountDecode)
		response.Status = codes.StatusError
		return response, nil
	}

	err = api.UnlockAccount(rpc.Account, rpc.DataPath, request.Passphrase)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
	}

	return response, nil
}

type requestSigninAccount struct {
	Name       string `mapstructure:"name"`
	Email      string `mapstructure:"email"`
	Passphrase string `mapstructure:"passphrase"`
}

// SigninAccount is the RPC method to sign in to an existing account
func SigninAccount(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	var request requestSigninAccount
	err := mapstructure.Decode(message.Payload, &request)
	if err != nil {
		rpc.Logger.Debug("Error decoding signin account request payload - ", err)
		response.Code = int(codes.ErrorSigninAccountDecode)
		response.Status = codes.StatusError
		return response, nil
	}

	newAccount, err := api.SigninAccount(rpc.DB, rpc.Logger, rpc.DataPath, request.Name, request.Email, request.Passphrase)
	if err == nil {
		rpc.Account = newAccount
	}
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
	}

	return response, nil
}

// SignoutAccount is the RPC method to sign out from the active account
func SignoutAccount(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	err := api.SignoutAccount(rpc.Account)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
	}
	rpc.Account = nil

	return response, nil
}

// LockAccount is the RPC method to lock the active account
func LockAccount(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	err := api.LockAccount(rpc.Account)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
	}

	return response, nil
}
