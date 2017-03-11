package rpc

import (
	"errors"

	"../account"
	"../crypto"
	pb "../proto"
	"golang.org/x/net/context"
)

// AccountState returns the accessible state of the account
func (rpc *Server) AccountState(ctx context.Context, request *pb.AccountStateRequest) (*pb.AccountStateResponse, error) {
	response := &pb.AccountStateResponse{
		SignedIn: false,
		Locked:   true,
	}
	if rpc.Account != nil {
		response.SignedIn = true
		response.Locked = rpc.Account.IsLocked()
	}
	return response, nil
}

// CreateAccount is the GRPC method to create a new account
func (rpc *Server) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// create account object
	newAccount := account.NewAccount(rpc.DB, rpc.Logger, request.Name)

	// create a new db file for the account
	err := newAccount.OpenAccountDb()
	if err != nil {
		return nil, err
	}
	// make this the active account
	rpc.Account = newAccount

	// create user object & attach it to the account
	user := account.NewUser(rpc.DB, rpc.Logger, request.Email)

	err = user.CreateKeys([]byte(request.Passphrase))
	if err != nil {
		return nil, err
	}

	// save user
	err = user.Save()
	if err != nil {
		return nil, err
	}
	newAccount.Users = append(newAccount.Users, user.Profile)
	newAccount.ActiveUser = user

	err = newAccount.Save()
	if err != nil {
		return nil, err
	}

	// creating the account automatically makes it the active account
	rpc.Account = newAccount

	response := &pb.CreateAccountResponse{
		Id: newAccount.ID.String(),
	}
	return response, nil
}

// UnlockAccount is the GRPC method to unlock the current account
func (rpc *Server) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.UnlockAccountResponse, error) {
	if rpc.Account == nil {
		return nil, errors.New("no active account")
	}

	if rpc.Account.ActiveUser == nil {
		return nil, errors.New("no active user")
	}

	// generate the derived key from the input passphrase and the stored salt
	key, err := crypto.DeriveKey([]byte(request.Passphrase), rpc.Account.ActiveUser.Salt)
	if err != nil {
		return nil, errors.New("error creating user key")
	}

	// encode the salt into the resulting key and store it in memory
	rpc.Account.ActiveUser.PassphraseKey = append(rpc.Account.ActiveUser.Salt, key[:]...)
	crypto.Zero(key[:])

	// since we never stored the original derived key
	// the only way we know if the key is valid is to try using it to open something
	_, err = crypto.Open(rpc.Account.ActiveUser.PassphraseKey, rpc.Account.ActiveUser.AccountKey)
	if err != nil {
		crypto.Zero(rpc.Account.ActiveUser.PassphraseKey)
		return nil, errors.New("invalid credentials")
	}

	err = rpc.Account.OpenAccountDb()
	if err != nil {
		return nil, errors.New("unable to open account db")
	}

	response := &pb.UnlockAccountResponse{}
	return response, nil
}

// SigninAccount is the GRPC method to sign in to an existing account
func (rpc *Server) SigninAccount(ctx context.Context, request *pb.SigninAccountRequest) (*pb.SigninAccountResponse, error) {
	// attempt to find the account (lookup)
	newAccount := account.NewAccount(rpc.DB, rpc.Logger, request.Name)
	err := newAccount.Lookup()
	if err != nil {
		return nil, errors.New("invalid account")
	}

	err = newAccount.OpenAccountDb()
	if err != nil {
		return nil, errors.New("unable to open account db")
	}

	// authenticate the user
	user := account.NewUser(rpc.DB, rpc.Logger, request.Email)
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
	newAccount.ActiveUser = user

	// load the account
	err = newAccount.Load()
	if err != nil {
		return nil, errors.New("unable to load account")
	}

	rpc.Account = newAccount

	// user should be signed in & account in an unlocked state at this point
	response := &pb.SigninAccountResponse{}

	return response, nil
}

// SignoutAccount is the GRPC method to sign out from the active account
func (rpc *Server) SignoutAccount(ctx context.Context, request *pb.SignoutAccountRequest) (*pb.SignoutAccountResponse, error) {
	if rpc.Account == nil {
		return nil, errors.New("no active account")
	}
	if rpc.Account.ActiveUser != nil {
		crypto.Zero(rpc.Account.ActiveUser.PassphraseKey)
	}
	rpc.Account.DB.Close()
	rpc.Account = nil

	response := &pb.SignoutAccountResponse{}
	return response, nil
}

// LockAccount is the GRPC method to lock the active account
func (rpc *Server) LockAccount(ctx context.Context, request *pb.LockAccountRequest) (*pb.LockAccountResponse, error) {
	crypto.Zero(rpc.Account.ActiveUser.PassphraseKey)
	rpc.Account.DB.Close()
	rpc.Account.DB = nil
	return nil, nil
}
