package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	pb "./proto"
	"github.com/boltdb/bolt"
	"golang.org/x/net/context"
	//	"google.golang.org/grpc/credentials"
)

// UIState returns the the UI state as saved by the master DB
func (backend *Backend) UIState(ctx context.Context, request *pb.UIStateRequest) (*pb.UIStateResponse, error) {
	response := &pb.UIStateResponse{}
	return response, nil
}

// AccountState returns the accessible state of the account
func (backend *Backend) AccountState(ctx context.Context, request *pb.AccountStateRequest) (*pb.AccountStateResponse, error) {
	response := &pb.AccountStateResponse{
		SignedIn: false,
		Locked:   true,
	}
	if backend.Account != nil {
		response.SignedIn = true
		response.Locked = backend.Account.IsLocked()
	}
	return response, nil
}

// OpenMasterDb opens the master database in the requested directory
func (backend *Backend) OpenMasterDb(ctx context.Context, request *pb.OpenMasterDbRequest) (*pb.OpenMasterDbResponse, error) {
	// This is the master index db
	// There are additional databases where actual notebook data is stored (one DB file per account)
	fileName := fmt.Sprint(filepath.Clean(request.Path), filepath.Separator, MasterDbFile)
	backend.Logger.Info("Opening master db file [", fileName, "]")
	var err error
	backend.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		backend.Logger.Error("Unable to open DB - ", err)
		return nil, err
	}
	response := &pb.OpenMasterDbResponse{}
	return response, nil
}

// CreateAccount is the GRPC method to create a new account
func (backend *Backend) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// create account object
	account := NewAccount(backend.DB, backend.Logger, request.Name)

	// create a new db file for the account
	err := account.OpenAccountDb()
	if err != nil {
		return nil, err
	}
	// make this the active account
	backend.Account = account

	// create user object & attach it to the account
	user := NewUser(account.DB, backend.Logger, request.Email)

	err = user.CreateKeys([]byte(request.Passphrase))
	if err != nil {
		return nil, err
	}

	// save user
	err = user.Save()
	if err != nil {
		return nil, err
	}
	account.Users = append(account.Users, user.Profile)
	account.ActiveUser = user

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

	err = backend.Account.OpenAccountDb()
	if err != nil {
		return nil, errors.New("unable to open account db")
	}

	response := &pb.UnlockAccountResponse{}
	return response, nil
}

// SigninAccount is the GRPC method to sign in to an existing account
func (backend *Backend) SigninAccount(ctx context.Context, request *pb.SigninAccountRequest) (*pb.SigninAccountResponse, error) {
	// attempt to find the account (lookup)
	account := NewAccount(backend.DB, backend.Logger, request.Name)
	err := account.Lookup()
	if err != nil {
		return nil, errors.New("invalid account")
	}

	err = account.OpenAccountDb()
	if err != nil {
		return nil, errors.New("unable to open account db")
	}

	// authenticate the user
	user := NewUser(account.DB, backend.Logger, request.Email)
	err = user.Lookup()
	if err != nil {
		return nil, errors.New("invalid user")
	}
	err = user.Load(request.Passphrase)
	if err != nil {
		// this error is probably because the passphrase was incorrect
		return nil, errors.New("unable to load user")
	}

	// connect the user to the account & make it the active user
	account.ActiveUser = user

	// load the account
	err = account.Load()
	if err != nil {
		return nil, errors.New("unable to load account")
	}

	// user should be signed in & account in an unlocked state at this point
	response := &pb.SigninAccountResponse{}

	return response, nil
}

// SignoutAccount is the GRPC method to sign out from the active account
func (backend *Backend) SignoutAccount(ctx context.Context, request *pb.SignoutAccountRequest) (*pb.SignoutAccountResponse, error) {
	if backend.Account == nil {
		return nil, errors.New("no active account")
	}
	if backend.Account.ActiveUser != nil {
		Zero(backend.Account.ActiveUser.PassphraseKey)
	}
	backend.Account.DB.Close()
	backend.Account = nil

	response := &pb.SignoutAccountResponse{}
	return response, nil
}

// LockAccount is the GRPC method to lock the active account
func (backend *Backend) LockAccount(ctx context.Context, request *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	Zero(backend.Account.ActiveUser.PassphraseKey)
	backend.Account.DB.Close()
	backend.Account.DB = nil
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
