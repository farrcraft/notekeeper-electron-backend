package handler

import (
	"../rpc"
)

// Handlers returns all available rpc handlers
func Handlers() map[string]rpc.Handler {
	handlers := make(map[string]rpc.Handler, 0)
	handlers["KeyExchange"] = KeyExchange

	handlers["MasterDb::open"] = OpenMasterDb

	handlers["Account::create"] = CreateAccount
	handlers["Account::unlock"] = UnlockAccount
	handlers["Account::signin"] = SigninAccount
	handlers["Account::signout"] = SignoutAccount
	handlers["Account::lock"] = LockAccount

	handlers["AccountState::get"] = GetAccountState

	handlers["UIState::load"] = LoadUIState
	handlers["UIState::save"] = SaveUIState

	handlers["User::shelves"] = GetUserShelves
	handlers["User::Shelf::create"] = CreateUserShelf
	handlers["User::Shelf::save"] = SaveUserShelf
	handlers["User::Shelf::delete"] = DeleteUserShelf

	handlers["Account::shelves"] = GetAccountShelves
	handlers["Account::Shelf::create"] = CreateAccountShelf
	handlers["Account::Shelf::save"] = SaveAccountShelf
	handlers["Account::Shelf::delete"] = DeleteAccountShelf

	handlers["User::collections"] = GetUserCollections
	handlers["User::Collection::create"] = CreateUserCollection
	handlers["User::Collection::save"] = SaveUserCollection
	handlers["User::Collection::delete"] = DeleteUserCollection

	handlers["Account::collections"] = GetAccountCollections
	handlers["Account::Collection::create"] = CreateAccountCollection
	handlers["Account::Collection::save"] = SaveAccountCollection
	handlers["Account::Collection::delete"] = DeleteAccountCollection

	handlers["User::tags"] = GetUserTags
	handlers["User::Tag::create"] = CreateUserTag
	handlers["User::Tag::save"] = SaveUserTag
	handlers["User::Tag::delete"] = DeleteUserTag

	handlers["Account::tags"] = GetAccountTags
	handlers["Account::Tag::create"] = CreateAccountTag
	handlers["Account::Tag::save"] = SaveAccountTag
	handlers["Account::Tag::delete"] = DeleteAccountTag

	handlers["User::notebooks"] = GetUserNotebooks
	handlers["User::Notebook::create"] = CreateUserNotebook
	handlers["User::Notebook::save"] = SaveUserNotebook
	handlers["User::Notebook::delete"] = DeleteUserNotebook

	handlers["Account::notebooks"] = GetAccountNotebooks
	handlers["Account::Notebook::create"] = CreateAccountNotebook
	handlers["Account::Notebook::save"] = SaveAccountNotebook
	handlers["Account::Notebook::delete"] = DeleteAccountNotebook

	handlers["User::notes"] = GetUserNotes
	handlers["User::Note::load"] = LoadUserNote
	handlers["User::Note::create"] = CreateUserNote
	handlers["User::Note::save"] = SaveUserNote
	handlers["User::Note::delete"] = DeleteUserNote

	handlers["Account::notes"] = GetAccountNotes
	handlers["Account::Note::load"] = LoadAccountNote
	handlers["Account::Note::create"] = CreateAccountNote
	handlers["Account::Note::save"] = SaveAccountNote
	handlers["Account::Note::delete"] = DeleteAccountNote

	return handlers
}
