package rpc

import (
	"../account"
	"../api"
	"../codes"
	"../db"
	messages "../proto"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// GetAccountState returns the accessible state of the account
func GetAccountState(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.AccountStateResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
		SignedIn: false,
		Locked:   true,
		Exists:   false,
	}

	if rpc.Account != nil {
		response.SignedIn = true
		response.Locked = rpc.Account.IsLocked()
		response.Exists = true
	} else {
		db := rpc.DBFactory.Find(db.TypeMaster, uuid.Nil)
		count := account.MapCount(db.DB)
		if count > 0 {
			response.Exists = true
		}
	}

	return response, nil
}

// CreateAccount is the RPC method to create a new account
func CreateAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	request := messages.CreateAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling create account request - ", err)
		response.Header.Code = int32(codes.ErrorCreateAccountDecode)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	// create the account
	newAccount, err := api.CreateAccount(rpc.DBFactory, rpc.Logger, request.Name, request.Email, request.Passphrase)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Status = code.Error()
		response.Header.Code = int32(code.Code)
		return response, nil
	}

	// make this the active account
	if err == nil {
		rpc.Account = newAccount
	}

	response.Id = newAccount.ID.String()

	return response, nil
}

// UnlockAccount is the RPC method to unlock the current account
func UnlockAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	request := messages.UnlockAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling unlock account request - ", err)
		response.Header.Code = int32(codes.ErrorUnlockAccountDecode)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	err = api.UnlockAccount(rpc.Account, request.Passphrase)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
	}

	return response, nil
}

// SigninAccount is the RPC method to sign in to an existing account
func SigninAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	request := messages.SigninAccountRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling signin account request - ", err)
		response.Header.Code = int32(codes.ErrorSigninAccountDecode)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	newAccount, err := api.SigninAccount(rpc.DBFactory, rpc.Logger, request.Name, request.Email, request.Passphrase)
	if err == nil {
		rpc.Account = newAccount
	}
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
	}

	return response, nil
}

// SignoutAccount is the RPC method to sign out from the active account
func SignoutAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	err := api.SignoutAccount(rpc.Account)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
	}
	rpc.Account = nil

	return response, nil
}

// LockAccount is the RPC method to lock the active account
func LockAccount(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	err := api.LockAccount(rpc.Account)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
	}

	return response, nil
}
