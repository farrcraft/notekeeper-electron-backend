package rpc

import (
	"../account"
	"../api"
	"../codes"
	pb "../proto"
	"golang.org/x/net/context"
)

// AccountState returns the accessible state of the account
func (rpc *Server) AccountState(ctx context.Context, request *pb.TokenRequest) (*pb.AccountStateResponse, error) {
	response := &pb.AccountStateResponse{
		SignedIn: false,
		Locked:   true,
		Exists:   false,
	}
	if rpc.Account != nil {
		response.SignedIn = true
		response.Locked = rpc.Account.IsLocked()
		response.Exists = true
	} else {
		count := account.MapCount(rpc.DB)
		if count > 0 {
			response.Exists = true
		}
	}
	return response, nil
}

// CreateAccount is the GRPC method to create a new account
func (rpc *Server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.IdResponse, error) {
	// create the account
	newAccount, err := api.CreateAccount(rpc.DB, rpc.Logger, rpc.DataPath, request.Name, request.Email, request.Passphrase)

	// make this the active account
	if err == nil {
		rpc.Account = newAccount
	}

	response := &pb.IdResponse{
		Id:     newAccount.ID.String(),
		Status: codes.StatusOK,
	}

	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int32(code.Code)
	}

	return response, nil
}

// UnlockAccount is the GRPC method to unlock the current account
func (rpc *Server) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.StatusResponse, error) {
	response := &pb.StatusResponse{
		Status: codes.StatusOK,
	}

	err := api.UnlockAccount(rpc.Account, rpc.DataPath, request.Passphrase)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int32(code.Code)
	}

	return response, nil
}

// SigninAccount is the GRPC method to sign in to an existing account
func (rpc *Server) SigninAccount(ctx context.Context, request *pb.SigninAccountRequest) (*pb.IdResponse, error) {
	newAccount, err := api.SigninAccount(rpc.DB, rpc.Logger, rpc.DataPath, request.Name, request.Email, request.Passphrase)
	if err == nil {
		rpc.Account = newAccount
	}

	// user should be signed in & account in an unlocked state at this point
	response := &pb.IdResponse{
		Status: codes.StatusOK,
	}

	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int32(code.Code)
	}

	return response, nil
}

// SignoutAccount is the GRPC method to sign out from the active account
func (rpc *Server) SignoutAccount(ctx context.Context, request *pb.IdRequest) (*pb.StatusResponse, error) {
	response := &pb.StatusResponse{
		Status: codes.StatusOK,
	}
	err := api.SignoutAccount(rpc.Account)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int32(code.Code)
	}
	rpc.Account = nil

	return response, nil
}

// LockAccount is the GRPC method to lock the active account
func (rpc *Server) LockAccount(ctx context.Context, request *pb.IdRequest) (*pb.StatusResponse, error) {
	response := &pb.StatusResponse{
		Status: codes.StatusOK,
	}
	err := api.LockAccount(rpc.Account)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int32(code.Code)
	}
	return response, nil
}
