package api

import (
	"../account"
	"../codes"
	"../crypto"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// CreateAccount creates a new account
func CreateAccount(db *bolt.DB, logger *logrus.Logger, dataPath string, name string, email string, passphrase string) (*account.Account, error) {
	// create account object
	newAccount := account.NewAccount(db, logger, name)

	// create a new db file for the account
	err := newAccount.OpenAccountDb(dataPath)
	if err != nil {
		return newAccount, err
	}

	// create user object & attach it to the account
	user := account.NewUser(db, logger, email)

	err = user.CreateKeys([]byte(passphrase))
	if err != nil {
		return newAccount, err
	}

	// save user
	err = user.Save()
	if err != nil {
		return newAccount, err
	}
	newAccount.Users = append(newAccount.Users, user.Profile)
	newAccount.ActiveUser = user

	err = newAccount.Save()
	if err != nil {
		return newAccount, err
	}
	return newAccount, nil
}

// UnlockAccount unlocks an account
func UnlockAccount(acct *account.Account, dataPath string, passphrase string) error {
	if acct == nil {
		code := codes.New(codes.ErrorUnlockActiveAccount)
		return code
	}

	if acct.ActiveUser == nil {
		code := codes.New(codes.ErrorUnlockActiveUser)
		return code
	}

	// generate the derived key from the input passphrase and the stored salt
	key, err := crypto.DeriveKey([]byte(passphrase), acct.ActiveUser.Salt)
	if err != nil {
		return err
	}

	// encode the salt into the resulting key and store it in memory
	acct.ActiveUser.PassphraseKey = key[:]

	// since we never stored the original derived key
	// the only way we know if the key is valid is to try using it to open something
	_, err = crypto.Open(acct.ActiveUser.PassphraseKey, acct.ActiveUser.AccountKey)
	if err != nil {
		crypto.Zero(acct.ActiveUser.PassphraseKey)
		return err
	}

	err = acct.OpenAccountDb(dataPath)
	if err != nil {
		return err
	}

	return nil
}

// SigninAccount signs in to an account
func SigninAccount(db *bolt.DB, logger *logrus.Logger, dataPath string, name string, email string, passphrase string) (*account.Account, error) {
	// attempt to find the account (lookup)
	newAccount := account.NewAccount(db, logger, name)
	err := newAccount.Lookup()
	if err != nil {
		return nil, err
	}

	err = newAccount.OpenAccountDb(dataPath)
	if err != nil {
		return nil, err
	}

	// authenticate the user
	user := account.NewUser(db, logger, email)
	err = user.Lookup()
	if err != nil {
		return nil, err
	}
	err = user.Load(passphrase)
	if err != nil {
		return nil, err
	}

	// connect the user to the account & make it the active user
	newAccount.ActiveUser = user

	// load the account
	err = newAccount.Load()
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}

// SignoutAccount signs out an account
func SignoutAccount(acct *account.Account) error {
	if acct == nil {
		code := codes.New(codes.ErrorSignoutActiveAccount)
		return code
	}
	if acct.ActiveUser == nil {
		code := codes.New(codes.ErrorSignoutActiveUser)
		return code
	}
	crypto.Zero(acct.ActiveUser.PassphraseKey)
	acct.DB.Close()
	return nil
}

// LockAccount locks an account
func LockAccount(acct *account.Account) error {
	if acct == nil {
		code := codes.New(codes.ErrorLockActiveAccount)
		return code
	}
	if acct.ActiveUser == nil {
		code := codes.New(codes.ErrorLockActiveUser)
		return code
	}
	crypto.Zero(acct.ActiveUser.PassphraseKey)
	if acct.DB == nil {
		code := codes.New(codes.ErrorLockActiveDb)
		return code
	}
	acct.DB.Close()
	acct.DB = nil
	return nil
}
