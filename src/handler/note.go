package handler

import (
	"../codes"
	"../note"
	messages "../proto"
	"../rpc"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// GetNotes is the RPC method to get a list of notes
func GetNotes(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.GetNotesResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetNotesRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling get notes request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope note.Scope
	var store note.StoreType
	if request.Scope == "account" {
		scope = note.ScopeAccount
	} else if request.Scope == "user" {
		scope = note.ScopeUser
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
		server.Logger.Debug("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Debug("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new note instance to act as a proxy
	n := note.New(nil, scope, store, server.DBFactory, server.Logger)
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

// LoadNote is the RPC method to get a note
func LoadNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.LoadNoteResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.LoadNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling load note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope note.Scope
	var store note.StoreType
	if request.Scope == "account" {
		scope = note.ScopeAccount
	} else if request.Scope == "user" {
		scope = note.ScopeUser
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
		server.Logger.Debug("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Debug("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new note instance to act as a proxy
	n := note.New(nil, scope, store, server.DBFactory, server.Logger)
	n.OwnerID = ownerID
	n.StoreID = storeID

	err = n.Load(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	response.Note = &messages.Note{
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

	return response, nil
}

// CreateNote is the RPC method to create a new note
func CreateNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope note.Scope
	var store note.StoreType
	if request.Scope == "account" {
		scope = note.ScopeAccount
	} else if request.Scope == "user" {
		scope = note.ScopeUser
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
		server.Logger.Debug("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Debug("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	n := note.New(t, scope, store, server.DBFactory, server.Logger)
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

// SaveNote is the RPC method to save an existing note
func SaveNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling save note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope note.Scope
	var store note.StoreType
	if request.Scope == "account" {
		scope = note.ScopeAccount
	} else if request.Scope == "user" {
		scope = note.ScopeUser
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
		server.Logger.Debug("Invalid note id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Debug("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	n := note.New(t, scope, store, server.DBFactory, server.Logger)
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

// DeleteNote is the RPC method to delete a note
func DeleteNote(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.DeleteNoteRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling load note request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope note.Scope
	var store note.StoreType
	if request.Scope == "account" {
		scope = note.ScopeAccount
	} else if request.Scope == "user" {
		scope = note.ScopeUser
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
		server.Logger.Debug("Invalid note id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid note owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	storeID, err := uuid.FromString(request.StoreId)
	if err != nil {
		server.Logger.Debug("Invalid note store id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	n := note.New(nil, scope, store, server.DBFactory, server.Logger)
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
