package main

import (
	pb "./proto"
	"golang.org/x/net/context"
	//	"google.golang.org/grpc/credentials"
	//	"golang.org/x/crypto/nacl/box"
)

// CreateAccount is the GRPC method to create a new account
func (backend *Backend) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	account := NewAccount(backend.DB, backend.Logger, request.Name)
	err := backend.CreateAccountDb(account)
	if err != nil {
		return nil, err
	}
	err = account.Save()
	if err != nil {
		return nil, err
	}
	response := &pb.CreateAccountResponse{
		Id: account.ID.String(),
	}
	return response, nil
}

// UnlockAccount is the GRPC method to unlock the current account
func (backend *Backend) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	return nil, nil
}

// SigninAccount is the GRPC method to sign in to an existing account
func (backend *Backend) SigninAccount(ctx context.Context, request *pb.SigninAccountRequest) (*pb.SigninAccountResponse, error) {
	return nil, nil
}

// CreateNotebook is the GRPC method to create a new notebook
func (backend *Backend) CreateNotebook(ctx context.Context, request *pb.CreateNotebookRequest) (*pb.CreateNotebookResponse, error) {
	notebook := NewNotebook(backend.Account)
	err := notebook.Save()
	if err != nil {
		return nil, err
	}
	response := &pb.CreateNotebookResponse{
		Status: "OK",
		Id:     notebook.ID.String(),
	}
	return response, nil
}
