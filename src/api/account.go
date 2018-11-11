package api

import (
	"../account"
	"../crypto"
	"../db"
	"../notebook"
	"../shelf"
	"../title"
	"../user"
)

// CreateAccount creates a new account
func (api *API) CreateAccount(name string, email string, passphrase string) (*account.Account, error) {
	// create account object
	newAccount, err := account.New(api.DBRegistry, api.Logger, name)
	if err != nil {
		return newAccount, err
	}

	// create a new db file for the account
	accountDBKey := db.Key{
		ID:   newAccount.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := api.DBRegistry.NewHandle(accountDBKey)
	if err != nil {
		return newAccount, err
	}

	// create user object & attach it to the account
	newUser, err := user.New(api.DBRegistry, api.Logger, newAccount.ID, email)
	if err != nil {
		return newAccount, err
	}

	err = newUser.CreateKeys([]byte(passphrase))
	if err != nil {
		return newAccount, err
	}

	userDBKey := db.Key{
		ID:   newUser.ID,
		Type: db.TypeUser,
	}
	userDBHandle, err := api.DBRegistry.NewHandle(userDBKey)
	if err != nil {
		return newAccount, err
	}

	// [FIXME] - save user mapping in user index in account db

	// save user
	err = newUser.Save()
	if err != nil {
		return newAccount, err
	}

	// [FIXME] - save user profile in account db
	// right now users are just stored as part of the account profile data
	newAccount.Users = append(newAccount.Users, newUser.Profile)
	newAccount.ActiveUser = newUser

	newAccount.EncryptedKey = newUser.AccountKey
	accountDBHandle.EncryptedKey = newUser.AccountKey
	userDBHandle.EncryptedKey = newUser.UserKey

	err = newAccount.Save()
	if err != nil {
		api.Logger.Debug("Error saving account - ", err)
		return newAccount, err
	}

	accountIndex := account.NewIndex(api.DBRegistry, api.Logger)
	err = accountIndex.Save(newAccount)
	if err != nil {
		api.Logger.Debug("Error saving account index - ", err)
		return newAccount, err
	}

	err = api.CreateAccountDefaults(newAccount, newUser)
	if err != nil {
		return newAccount, err
	}

	return newAccount, nil
}

// CreateAccountDefaults creates default objects and settings for an account
func (api *API) CreateAccountDefaults(acct *account.Account, currentUser *user.User) error {
	// Create the account-scoped default shelf 'My Shelf'
	defaultShelfTitle := title.New("My Shelf")

	accountShelf, err := shelf.New(defaultShelfTitle, shelf.ScopeAccount, api.DBRegistry, api.Logger)
	if err != nil {
		return err
	}

	accountShelf.OwnerID = acct.ID
	accountShelf.Default = true

	shelfDBKey := db.Key{
		ID:   accountShelf.ID,
		Type: db.TypeShelf,
	}
	accountShelfDBHandle, err := api.DBRegistry.NewHandle(shelfDBKey)
	if err != nil {
		return err
	}

	accountShelf.EncryptedKey, err = acct.CreateEncryptedKey()
	if err != nil {
		return err
	}
	accountShelfDBHandle.EncryptedKey = accountShelf.EncryptedKey

	// shelf metadata will be saved in an index (user or account db)
	accountShelfIndex := shelf.NewIndex(shelf.ScopeAccount, acct.ID, api.DBRegistry, api.Logger)
	unsealedAccountKey, err := acct.UnsealKey(account.TypePassphrase, acct.EncryptedKey)
	if err != nil {
		return err
	}
	err = accountShelfIndex.Save(accountShelf, unsealedAccountKey)
	if err != nil {
		api.Logger.Debug("could not create default account shelf")
		return err
	}

	// Create the user-scoped default shelf 'My Shelf'
	userShelf, err := shelf.New(defaultShelfTitle, shelf.ScopeUser, api.DBRegistry, api.Logger)
	if err != nil {
		return err
	}

	userShelf.OwnerID = currentUser.ID
	userShelf.Default = true

	shelfDBKey.ID = userShelf.ID
	userShelfDBHandle, err := api.DBRegistry.NewHandle(shelfDBKey)
	if err != nil {
		return err
	}

	userShelf.EncryptedKey, err = currentUser.CreateEncryptedKey(user.TypeUser)
	if err != nil {
		return err
	}
	userShelfDBHandle.EncryptedKey = userShelf.EncryptedKey

	userShelfIndex := shelf.NewIndex(shelf.ScopeUser, currentUser.ID, api.DBRegistry, api.Logger)
	unsealedUserKey, err := acct.UnsealKey(account.TypePassphrase, acct.ActiveUser.UserKey)
	if err != nil {
		api.Logger.Debug("could not unseal user key")
		return err
	}
	err = userShelfIndex.Save(userShelf, unsealedUserKey)
	if err != nil {
		api.Logger.Debug("could not create default user shelf")
		return err
	}

	// Create the account-scoped default notebook 'My Notebook' inside the account-scoped default shelf
	defaultNotebookTitle := title.New("My Notebook")

	accountNotebook, err := notebook.New(defaultNotebookTitle, notebook.ScopeAccount, notebook.ContainerTypeShelf, api.DBRegistry, api.Logger)
	if err != nil {
		return err
	}

	accountNotebook.OwnerID = acct.ID
	accountNotebook.ContainerID = accountShelf.ID
	accountNotebook.Default = true
	accountShelfKey, err := acct.UnsealKey(account.TypeAccount, accountShelfDBHandle.EncryptedKey)
	if err != nil {
		api.Logger.Debug("could not unseal default account shelf key")
		return err
	}
	err = accountNotebook.Save(accountShelfKey)
	if err != nil {
		api.Logger.Debug("could not create default account notebook")
		return err
	}

	// Create the user-scoped default notebook 'My Notebook' inside the user-scoped default shelf
	userNotebook, err := notebook.New(defaultNotebookTitle, notebook.ScopeUser, notebook.ContainerTypeShelf, api.DBRegistry, api.Logger)
	if err != nil {
		return err
	}

	userNotebook.OwnerID = currentUser.ID
	userNotebook.ContainerID = userShelf.ID
	userNotebook.Default = true
	userShelfKey, err := currentUser.UnsealKey(user.TypeUser, userShelfDBHandle.EncryptedKey)
	if err != nil {
		api.Logger.Debug("could not unseal default user shelf key")
		return err
	}
	err = userNotebook.Save(userShelfKey)
	if err != nil {
		api.Logger.Debug("could not create default user notebook")
		return err
	}

	// Create the account-scoped special 'Trash' shelf
	trashShelfTitle := title.New("Trash")

	accountTrashShelf, err := shelf.New(trashShelfTitle, shelf.ScopeAccount, api.DBRegistry, api.Logger)
	if err != nil {
		return err
	}

	accountTrashShelf.OwnerID = acct.ID
	accountTrashShelf.Trash = true

	shelfDBKey.ID = accountTrashShelf.ID
	accountTrashShelfDBHandle, err := api.DBRegistry.NewHandle(shelfDBKey)
	if err != nil {
		return err
	}

	accountTrashShelf.EncryptedKey, err = acct.CreateEncryptedKey()
	if err != nil {
		return err
	}
	accountTrashShelfDBHandle.EncryptedKey = accountTrashShelf.EncryptedKey

	err = accountShelfIndex.Save(accountTrashShelf, unsealedAccountKey)
	if err != nil {
		api.Logger.Debug("could not create account trash shelf")
		return err
	}

	// Create the user-scoped special 'Trash' shelf
	userTrashShelf, err := shelf.New(trashShelfTitle, shelf.ScopeUser, api.DBRegistry, api.Logger)
	if err != nil {
		return err
	}

	userTrashShelf.OwnerID = currentUser.ID
	userTrashShelf.Trash = true

	shelfDBKey.ID = userTrashShelf.ID
	userTrashShelfDBHandle, err := api.DBRegistry.NewHandle(shelfDBKey)
	if err != nil {
		return err
	}

	userTrashShelf.EncryptedKey, err = acct.CreateEncryptedKey()
	if err != nil {
		return err
	}
	userTrashShelfDBHandle.EncryptedKey = accountTrashShelf.EncryptedKey

	err = userShelfIndex.Save(userTrashShelf, unsealedUserKey)
	if err != nil {
		api.Logger.Debug("could not create user trash shelf")
		return err
	}

	return nil
}

// SigninAccount signs in to an account
func (api *API) SigninAccount(name string, email string, passphrase string) (*account.Account, error) {
	// attempt to find the account (lookup)
	newAccount, err := account.New(api.DBRegistry, api.Logger, name)
	if err != nil {
		return nil, err
	}

	accountIndex := account.NewIndex(api.DBRegistry, api.Logger)
	err = accountIndex.Lookup(newAccount)
	if err != nil {
		return nil, err
	}

	accountDBKey := db.Key{
		ID:   newAccount.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := api.DBRegistry.NewHandle(accountDBKey)
	if err != nil {
		return nil, err
	}

	// authenticate the user
	newUser, err := user.New(api.DBRegistry, api.Logger, newAccount.ID, email)
	if err != nil {
		return nil, err
	}

	// resolve the user id from the user map in the account db
	userIndex := user.NewIndex(newAccount.ID, api.DBRegistry, api.Logger)
	err = userIndex.Lookup(newUser)
	if err != nil {
		api.DBRegistry.CloseAccountDBs()
		return nil, err
	}
	// load the user from the user db
	userDBKey := db.Key{
		ID:   newUser.ID,
		Type: db.TypeUser,
	}
	_, err = api.DBRegistry.NewHandle(userDBKey)
	if err != nil {
		return nil, err
	}
	err = newUser.Load(passphrase)
	if err != nil {
		api.DBRegistry.CloseAccountDBs()
		return nil, err
	}

	// connect the user to the account & make it the active user
	newAccount.ActiveUser = newUser

	// set remaining encryption keys
	newAccount.EncryptedKey = newUser.AccountKey
	accountDBHandle.EncryptedKey = newUser.AccountKey

	// load the account
	err = newAccount.Load()
	if err != nil {
		api.DBRegistry.CloseAccountDBs()
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
	if acct.DBRegistry == nil {
		api.Logger.Debug("signout missing db registry")
		return nil
	}
	acct.DBRegistry.CloseAccountDBs()
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
	acct.ActiveUser.UserKey = []byte{}
	acct.ActiveUser.AccountKey = []byte{}
	acct.EncryptedKey = []byte{}
	//crypto.Zero(acct.ActiveUser.PassphraseKey)
	if acct.DBRegistry == nil {
		api.Logger.Debug("lock account missing db factory")
		return nil
	}
	acct.DBRegistry.CloseAccountDBs()
	return nil
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

	accountDBKey := db.Key{
		ID:   acct.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := api.DBRegistry.NewHandle(accountDBKey)
	if err != nil {
		api.Logger.Debug("unlock could not open account db")
		return err
	}

	// load the user db
	userDBKey := db.Key{
		ID:   acct.ActiveUser.ID,
		Type: db.TypeUser,
	}
	_, err = api.DBRegistry.NewHandle(userDBKey)
	if err != nil {
		return err
	}

	err = acct.ActiveUser.Load(passphrase)
	if err != nil {
		api.DBRegistry.CloseAccountDBs()
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

	accountDBHandle.EncryptedKey = acct.ActiveUser.AccountKey

	return nil
}
