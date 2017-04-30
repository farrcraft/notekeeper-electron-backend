package api

import (
	"../account"
	"../crypto"
	"../user"
)

// CreateAccount creates a new account
func (api *API) CreateAccount(name string, email string, passphrase string) (*account.Account, error) {
	// create account object
	newAccount := account.New(api.DBFactory, api.Logger, name)

	// create a new db file for the account
	err := newAccount.OpenAccountDb()
	if err != nil {
		return newAccount, err
	}

	// create user object & attach it to the account
	user := user.New(api.DBFactory, api.Logger, newAccount.ID, email)

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
func (api *API) UnlockAccount(acct *account.Account, passphrase string) error {
	if acct == nil {
		api.Logger.Debug("unlock missing account")
		return nil
	}

	if acct.ActiveUser == nil {
		api.Logger.Debug("unlock missing user")
		return nil
	}

	err := acct.OpenAccountDb()
	if err != nil {
		api.Logger.Debug("unlock could not open account db")
		return err
	}

	err = acct.ActiveUser.Lookup()
	if err != nil {
		api.Logger.Debug("unlock could not lookup user")
		return err
	}

	// generate the derived key from the input passphrase and the stored salt
	c := crypto.New(api.Logger)
	key, err := c.DeriveKey([]byte(passphrase), acct.ActiveUser.Salt)
	if err != nil {
		return err
	}

	// encode the salt into the resulting key and store it in memory
	acct.ActiveUser.PassphraseKey = key[:]

	// since we never stored the original derived key
	// the only way we know if the key is valid is to try using it to open something
	_, err = c.Open(acct.ActiveUser.PassphraseKey, acct.ActiveUser.AccountKey)
	if err != nil {
		crypto.Zero(acct.ActiveUser.PassphraseKey)
		return err
	}

	err = acct.OpenAccountDb()
	if err != nil {
		return err
	}

	return nil
}

// SigninAccount signs in to an account
func (api *API) SigninAccount(name string, email string, passphrase string) (*account.Account, error) {
	// attempt to find the account (lookup)
	newAccount := account.New(api.DBFactory, api.Logger, name)
	err := newAccount.Lookup()
	if err != nil {
		return nil, err
	}

	err = newAccount.OpenAccountDb()
	if err != nil {
		return nil, err
	}

	// authenticate the user
	user := user.New(api.DBFactory, api.Logger, newAccount.ID, email)
	// resolve the user id from the user map in the account db
	err = user.Lookup()
	if err != nil {
		api.DBFactory.CloseAccountDBs()
		return nil, err
	}
	// load the user from the user db
	err = user.Load(passphrase)
	if err != nil {
		api.DBFactory.CloseAccountDBs()
		return nil, err
	}

	// connect the user to the account & make it the active user
	newAccount.ActiveUser = user

	// load the account
	err = newAccount.Load()
	if err != nil {
		api.DBFactory.CloseAccountDBs()
		return nil, err
	}
	return newAccount, nil
}

// SignoutAccount signs out an account
func (api *API) SignoutAccount(acct *account.Account) error {
	if acct == nil {
		api.Logger.Debug("signout missing account")
		return nil
	}
	if acct.ActiveUser == nil {
		api.Logger.Debug("signout missing user")
		return nil
	}
	crypto.Zero(acct.ActiveUser.PassphraseKey)
	if acct.DBFactory == nil {
		api.Logger.Debug("signout missing db factory")
		return nil
	}
	acct.DBFactory.CloseAccountDBs()
	return nil
}

// LockAccount locks an account
func (api *API) LockAccount(acct *account.Account) error {
	if acct == nil {
		api.Logger.Debug("lock account missing account")
		return nil
	}
	if acct.ActiveUser == nil {
		api.Logger.Debug("lock account missing user")
		return nil
	}
	acct.ActiveUser.PassphraseKey = []byte{}
	//crypto.Zero(acct.ActiveUser.PassphraseKey)
	if acct.DBFactory == nil {
		api.Logger.Debug("lock account missing db factory")
		return nil
	}
	acct.DBFactory.CloseAccountDBs()
	return nil
}
