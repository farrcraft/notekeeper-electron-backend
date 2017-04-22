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

	handlers["Shelf::loadAll"] = GetShelves
	handlers["Shelf::create"] = CreateShelf
	handlers["Shelf::save"] = SaveShelf
	handlers["Shelf::delete"] = DeleteShelf

	handlers["Collection::loadAll"] = GetCollections
	handlers["Collection::create"] = CreateCollection
	handlers["Collection::save"] = SaveCollection
	handlers["Collection::delete"] = DeleteCollection

	/*
		handlers["Tag::loadAll"] = GetTags
		handlers["Tag::create"] = CreateTag
		handlers["Tag::save"] = SaveTag
		handlers["Tag::delete"] = DeleteTag

		handlers["Notebook::loadAll"] = GetNotebooks
		handlers["Notebook::create"] = CreateNotebook
		handlers["Notebook::save"] = SaveNotebook
		handlers["Notebook::delete"] = DeleteNotebook

		handlers["Note::loadAll"] = GetNotes
		handlers["Note::load"] = LoadNote
		handlers["Note::create"] = CreateNote
		handlers["Note::save"] = SaveNote
		handlers["Note::delete"] = DeleteNote
	*/

	return handlers
}
