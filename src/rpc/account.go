package rpc

import (
	"errors"

	"../account"
	"../crypto"
	pb "../proto"
	"golang.org/x/net/context"
)

// AccountState returns the accessible state of the account
func (rpc *Server) AccountState(ctx context.Context, request *pb.TokenRequest) (*pb.AccountStateResponse, error) {
	rpc.Logger.Debug("rpc - getting account state")
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
	rpc.Logger.Debug("rpc - creating account")
	// create account object
	newAccount := account.NewAccount(rpc.DB, rpc.Logger, request.Name)

	// create a new db file for the account
	err := newAccount.OpenAccountDb(rpc.DataPath)
	if err != nil {
		rpc.Logger.Debug("rpc - create account error - unable to open account db")
		return nil, errors.New("unable to open account db")
	}
	// make this the active account
	rpc.Account = newAccount

	// create user object & attach it to the account
	user := account.NewUser(rpc.DB, rpc.Logger, request.Email)

	err = user.CreateKeys([]byte(request.Passphrase))
	if err != nil {
		rpc.Logger.Debug("rpc - create account error - unable to create keys")
		return nil, errors.New("unable to create keys")
	}

	// save user
	err = user.Save()
	if err != nil {
		rpc.Logger.Debug("rpc - create account error - error saving user")
		return nil, errors.New("error saving user")
	}
	newAccount.Users = append(newAccount.Users, user.Profile)
	newAccount.ActiveUser = user

	err = newAccount.Save()
	if err != nil {
		rpc.Logger.Debug("rpc - create account error - error saving account")
		return nil, errors.New("error saving account")
	}

	// creating the account automatically makes it the active account
	rpc.Account = newAccount

	response := &pb.IdResponse{
		Id: newAccount.ID.String(),
	}
	return response, nil
}

// UnlockAccount is the GRPC method to unlock the current account
func (rpc *Server) UnlockAccount(ctx context.Context, request *pb.UnlockAccountRequest) (*pb.StatusResponse, error) {
	rpc.Logger.Debug("rpc - unlocking account")
	if rpc.Account == nil {
		rpc.Logger.Debug("rpc - unlock account error - no active account")
		return nil, errors.New("no active account")
	}

	if rpc.Account.ActiveUser == nil {
		rpc.Logger.Debug("rpc - unlock account error - no active user")
		return nil, errors.New("no active user")
	}

	// generate the derived key from the input passphrase and the stored salt
	key, err := crypto.DeriveKey([]byte(request.Passphrase), rpc.Account.ActiveUser.Salt)
	if err != nil {
		rpc.Logger.Debug("rpc - unlock account error - error creating user key")
		return nil, errors.New("error creating user key")
	}

	// encode the salt into the resulting key and store it in memory
	rpc.Account.ActiveUser.PassphraseKey = append(rpc.Account.ActiveUser.Salt, key[:]...)
	crypto.Zero(key[:])

	// since we never stored the original derived key
	// the only way we know if the key is valid is to try using it to open something
	_, err = crypto.Open(rpc.Account.ActiveUser.PassphraseKey, rpc.Account.ActiveUser.AccountKey)
	if err != nil {
		rpc.Logger.Debug("rpc - unlock account error - invalid credentials")
		crypto.Zero(rpc.Account.ActiveUser.PassphraseKey)
		return nil, errors.New("invalid credentials")
	}

	err = rpc.Account.OpenAccountDb(rpc.DataPath)
	if err != nil {
		rpc.Logger.Debug("rpc - unlock account error - unable to open account db")
		return nil, errors.New("unable to open account db")
	}

	response := &pb.StatusResponse{}
	return response, nil
}

// SigninAccount is the GRPC method to sign in to an existing account
func (rpc *Server) SigninAccount(ctx context.Context, request *pb.SigninAccountRequest) (*pb.IdResponse, error) {
	rpc.Logger.Debug("rpc - signing into account")
	// attempt to find the account (lookup)
	newAccount := account.NewAccount(rpc.DB, rpc.Logger, request.Name)
	err := newAccount.Lookup()
	if err != nil {
		rpc.Logger.Debug("rpc - sign in error - invalid account")
		return nil, errors.New("invalid account")
	}

	err = newAccount.OpenAccountDb(rpc.DataPath)
	if err != nil {
		rpc.Logger.Debug("rpc - sign in error - unable to open account db")
		return nil, errors.New("unable to open account db")
	}

	// authenticate the user
	user := account.NewUser(rpc.DB, rpc.Logger, request.Email)
	err = user.Lookup()
	if err != nil {
		rpc.Logger.Debug("rpc - sign in error - invalid user")
		return nil, errors.New("invalid user")
	}
	err = user.Load(request.Passphrase)
	if err != nil {
		rpc.Logger.Debug("rpc - sign in error - unable to load user")
		// this error is probably because the passphrase was incorrect
		return nil, errors.New("unable to load user")
	}

	// connect the user to the account & make it the active user
	newAccount.ActiveUser = user

	// load the account
	err = newAccount.Load()
	if err != nil {
		rpc.Logger.Debug("rpc -sign in error - unable to load account")
		return nil, errors.New("unable to load account")
	}

	rpc.Account = newAccount

	// user should be signed in & account in an unlocked state at this point
	response := &pb.IdResponse{}

	return response, nil
}

// SignoutAccount is the GRPC method to sign out from the active account
func (rpc *Server) SignoutAccount(ctx context.Context, request *pb.IdRequest) (*pb.StatusResponse, error) {
	rpc.Logger.Debug("rpc - signing out of account")
	if rpc.Account == nil {
		rpc.Logger.Debug("rpc - signout error - no active account")
		return nil, errors.New("no active account")
	}
	if rpc.Account.ActiveUser == nil {
		rpc.Logger.Debug("rpc - signout error - no active user")
		return nil, errors.New("no active user")
	}
	crypto.Zero(rpc.Account.ActiveUser.PassphraseKey)
	rpc.Account.DB.Close()
	rpc.Account = nil

	response := &pb.StatusResponse{}
	return response, nil
}

// LockAccount is the GRPC method to lock the active account
func (rpc *Server) LockAccount(ctx context.Context, request *pb.IdRequest) (*pb.StatusResponse, error) {
	rpc.Logger.Debug("rpc - locking account")
	if rpc.Account == nil {
		rpc.Logger.Debug("rpc - lock account error - no active account")
		return nil, errors.New("no active account")
	}
	if rpc.Account.ActiveUser == nil {
		rpc.Logger.Debug("rpc - lock account error - no active user")
		return nil, errors.New("no active user")
	}
	crypto.Zero(rpc.Account.ActiveUser.PassphraseKey)
	rpc.Account.DB.Close()
	rpc.Account.DB = nil
	response := &pb.StatusResponse{}
	return response, nil
}
