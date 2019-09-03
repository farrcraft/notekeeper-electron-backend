package handler

import (
	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/note"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

func getNotes(server *rpc.Server, message []byte, scope string) (proto.Message, error) {
	response := &messages.GetNotesResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetNotesRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling get notes request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var noteScope note.Scope
	var store note.StoreType
	if scope == "account" {
		noteScope = note.ScopeAccount
	} else if scope == "user" {
		noteScope = note.ScopeUser
	} else {
		return response, nil
	}

	if request.Store == "collection" {
		store = note.StoreTypeCollection
	} else if request.Store == "shelf" {
		store = note.StoreTypeShelf
	} else {
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Warn("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new note instance to act as a proxy
	n, err := note.New(nil, noteScope, store, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating note - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	n.OwnerID = ownerID
	n.StoreID = storeID

	notes, err := n.LoadAll(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	for _, n := range notes {
		m := &messages.Note{
			Id:      n.ID.String(),
			Scope:   request.Scope,
			Store:   request.Store,
			OwnerId: request.OwnerId,
			StoreId: request.StoreId,
			Name:    rpc.TitleToMessage(n.Title),
			Locked:  n.Locked,
			Created: rpc.TimeToMessage(n.Created),
			Updated: rpc.TimeToMessage(n.Updated),
		}
		response.Notes = append(response.Notes, m)
	}

	return response, nil
}

func loadNote(server *rpc.Server, message []byte, scope string) (proto.Message, error) {
	response := &messages.LoadNoteResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.LoadNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling load note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var noteScope note.Scope
	var store note.StoreType
	if scope == "account" {
		noteScope = note.ScopeAccount
	} else if scope == "user" {
		noteScope = note.ScopeUser
	} else {
		return response, nil
	}

	if request.Store == "collection" {
		store = note.StoreTypeCollection
	} else if request.Store == "shelf" {
		store = note.StoreTypeShelf
	} else {
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Warn("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new note instance to act as a proxy
	n, err := note.New(nil, noteScope, store, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating note - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	n.OwnerID = ownerID
	n.StoreID = storeID

	err = n.Load(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	response.Note = &messages.Note{
		Id:      n.ID.String(),
		Store:   request.Store,
		OwnerId: request.OwnerId,
		StoreId: request.StoreId,
		Name:    rpc.TitleToMessage(n.Title),
		Locked:  n.Locked,
		Created: rpc.TimeToMessage(n.Created),
		Updated: rpc.TimeToMessage(n.Updated),
	}

	return response, nil
}

func createNote(server *rpc.Server, message []byte, scope string) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling create note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var noteScope note.Scope
	var store note.StoreType
	if scope == "account" {
		noteScope = note.ScopeAccount
	} else if scope == "user" {
		noteScope = note.ScopeUser
	} else {
		return response, nil
	}

	if request.Store == "collection" {
		store = note.StoreTypeCollection
	} else if request.Store == "shelf" {
		store = note.StoreTypeShelf
	} else {
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Warn("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	n, err := note.New(t, noteScope, store, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating note - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	n.OwnerID = ownerID
	n.StoreID = storeID

	err = n.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	response.Id = n.ID.String()

	return response, nil
}

func saveNote(server *rpc.Server, message []byte, scope string) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling save note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var noteScope note.Scope
	var store note.StoreType
	if scope == "account" {
		noteScope = note.ScopeAccount
	} else if scope == "user" {
		noteScope = note.ScopeUser
	} else {
		return response, nil
	}

	if request.Store == "collection" {
		store = note.StoreTypeCollection
	} else if request.Store == "shelf" {
		store = note.StoreTypeShelf
	} else {
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Warn("Invalid note id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Warn("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	n, err := note.New(t, noteScope, store, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating note - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	n.ID = id
	n.OwnerID = ownerID
	n.StoreID = storeID

	err = n.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}
	return response, nil
}

func deleteNote(server *rpc.Server, message []byte, scope string) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.DeleteNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling load note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var noteScope note.Scope
	var store note.StoreType
	if scope == "account" {
		noteScope = note.ScopeAccount
	} else if scope == "user" {
		noteScope = note.ScopeUser
	} else {
		return response, nil
	}

	if request.Store == "collection" {
		store = note.StoreTypeCollection
	} else if request.Store == "shelf" {
		store = note.StoreTypeShelf
	} else {
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Warn("Invalid note id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Warn("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	n, err := note.New(nil, noteScope, store, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating note - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	n.ID = id
	n.OwnerID = ownerID
	n.StoreID = storeID

	err = n.Delete(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}
	return response, nil
}
