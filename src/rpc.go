package main

import (
	"errors"

	pb "./proto"
	"golang.org/x/net/context"
	//	"google.golang.org/grpc/credentials"
)

// CreateAccount is the GRPC method to create a new account
func (backend *Backend) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// create account object
	account := NewAccount(backend.DB, backend.Logger, request.Name)

	// create a new db file for the account
	err := backend.CreateAccountDb(account)
	if err != nil {
		return nil, err
	}

	// create user object & attach it to the account
	user := NewUser(backend.DB, backend.Logger, request.Email)
	account.Users = append(account.Users, user)
	account.ActiveUser = user

	// generate account-level encryption key
	accountKey, err := GenerateKey()
	if err != nil {
		return nil, err
	}
	// derive key from passphrase
	var key = new([KeySize]byte)
	key, user.Salt, err = DeriveKeyAndSalt([]byte(request.Passphrase))
	if err != nil {
		return nil, err
	}
	slicedKey := key[:]
	user.PassphraseKey = append(user.Salt, slicedKey...)
	user.AccountKey, err = Seal(user.PassphraseKey, accountKey[:])
	if err != nil {
		return nil, err
	}
	Zero(accountKey[:])
	Zero(slicedKey)

	// [FIXME] - save user

	err = account.Save()
	if err != nil {
		return nil, err
	}

	// creating the account automatically makes it the active account
	backend.Account = account

	response := &pb.CreateAccountResponse{
		Id: account.ID.String(),
	}
	return response, nil
}

// UnlockAccount is the GRPC method to unlock the current account
func (backend *Backend) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	if backend.Account == nil {
		return nil, errors.New("no active account")
	}

	if backend.Account.ActiveUser == nil {
		return nil, errors.New("no active user")
	}

	// generate the derived key from the input passphrase and the stored salt
	key, err := DeriveKey([]byte(request.Passphrase), backend.Account.ActiveUser.Salt)
	if err != nil {
		return nil, errors.New("error creating user key")
	}

	// encode the salt into the resulting key and store it in memory
	backend.Account.ActiveUser.PassphraseKey = append(backend.Account.ActiveUser.Salt, key[:]...)
	Zero(key[:])

	// since we never stored the original derived key
	// the only way we know if the key is valid is to try using it to open something
	_, err = Open(backend.Account.ActiveUser.PassphraseKey, backend.Account.ActiveUser.AccountKey)
	if err != nil {
		Zero(backend.Account.ActiveUser.PassphraseKey)
		return nil, errors.New("invalid credentials")
	}

	response := &pb.UnlockAccountResponse{}
	return response, nil
}

// SigninAccount is the GRPC method to sign in to an existing account
func (backend *Backend) SigninAccount(ctx context.Context, request *pb.SigninAccountRequest) (*pb.SigninAccountResponse, error) {
	// attempt to find the account (lookup)
	// load the account
	// authenticate the user
	return nil, nil
}

// SignoutAccount is the GRPC method to sign out from the active account
func (backend *Backend) SignoutAccount(ctx context.Context, request *pb.SignoutAccountRequest) (*pb.SignoutAccountResponse, error) {
	if backend.Account == nil {
		return nil, errors.New("no active account")
	}
	if backend.Account.ActiveUser != nil {
		Zero(backend.Account.ActiveUser.PassphraseKey)
	}
	backend.Account = nil

	response := &pb.SignoutAccountResponse{}
	return response, nil
}

// LockAccount is the GRPC method to lock the active account
func (backend *Backend) LockAccount(ctx context.Context, request *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	Zero(backend.Account.ActiveUser.PassphraseKey)
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
