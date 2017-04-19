package rpc

// RegisterHandlers registers all of the RPC handlers
func (rpc *Server) RegisterHandlers() {
	rpc.Handlers["KeyExchange"] = KeyExchange

	rpc.Handlers["MasterDb::open"] = OpenMasterDb

	rpc.Handlers["Account::create"] = CreateAccount
	rpc.Handlers["Account::unlock"] = UnlockAccount
	rpc.Handlers["Account::signin"] = SigninAccount
	rpc.Handlers["Account::signout"] = SignoutAccount
	rpc.Handlers["Account::lock"] = LockAccount

	rpc.Handlers["AccountState::get"] = GetAccountState

	rpc.Handlers["UIState::load"] = LoadUIState
	rpc.Handlers["UIState::save"] = SaveUIState

	rpc.Handlers["Shelf::loadAll"] = GetShelves
	rpc.Handlers["Shelf::create"] = CreateShelf
	rpc.Handlers["Shelf::save"] = SaveShelf
	rpc.Handlers["Shelf::delete"] = DeleteShelf
	/*
		rpc.Handlers["Collection::loadAll"] = GetCollections
		rpc.Handlers["Collection::create"] = CreateCollection
		rpc.Handlers["Collection::save"] = SaveCollection
		rpc.Handlers["Collection::delete"] = DeleteCollection

		rpc.Handlers["Tag::loadAll"] = GetTags
		rpc.Handlers["Tag::create"] = CreateTag
		rpc.Handlers["Tag::save"] = SaveTag
		rpc.Handlers["Tag::delete"] = DeleteTag

		rpc.Handlers["Notebook::loadAll"] = GetNotebooks
		rpc.Handlers["Notebook::create"] = CreateNotebook
		rpc.Handlers["Notebook::save"] = SaveNotebook
		rpc.Handlers["Notebook::delete"] = DeleteNotebook

		rpc.Handlers["Note::loadAll"] = GetNotes
		rpc.Handlers["Note::load"] = LoadNote
		rpc.Handlers["Note::create"] = CreateNote
		rpc.Handlers["Note::save"] = SaveNote
		rpc.Handlers["Note::delete"] = DeleteNote
	*/
}
