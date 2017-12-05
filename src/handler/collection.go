package handler

import (
	"../codes"
	"../collection"
	messages "../proto"
	"../rpc"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// GetCollections gets a list of collections
func GetCollections(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.GetCollectionsResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetCollectionsRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling get collections request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	shelfID, err := uuid.FromString(request.ShelfId)
	if err != nil {
		server.Logger.Debug("Invalid shelf id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	var ownerID uuid.UUID
	var scope collection.Scope
	if request.Scope == "account" {
		ownerID = server.Account.ID
		scope = collection.ScopeAccount
	} else if request.Scope == "user" {
		ownerID = server.Account.ActiveUser.ID
		scope = collection.ScopeUser
	} else {
		return response, nil
	}

	// create a new collection instance to act as a proxy
	index := collection.NewIndex(scope, server.DBRegistry, server.Logger)
	index.ShelfID = shelfID
	index.OwnerID = ownerID

	err = index.LoadAll(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	for _, c := range index.Collections {
		m := &messages.Collection{
			Id:      c.ID.String(),
			ShelfId: shelfID.String(),
			Name:    rpc.TitleToMessage(c.Title),
			Locked:  c.Locked,
			Created: rpc.TimeToMessage(c.Created),
			Updated: rpc.TimeToMessage(c.Updated),
		}
		response.Collections = append(response.Collections, m)
	}

	return response, nil
}

// CreateCollection saves a new collection
func CreateCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateCollectionRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create collection request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	shelfID, err := uuid.FromString(request.ShelfId)
	if err != nil {
		server.Logger.Debug("Invalid shelf id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	var ownerID uuid.UUID
	var scope collection.Scope
	if request.Scope == "account" {
		ownerID = server.Account.ID
		scope = collection.ScopeAccount
	} else if request.Scope == "user" {
		ownerID = server.Account.ActiveUser.ID
		scope = collection.ScopeUser
	} else {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	c := collection.New(t, scope, server.DBRegistry, server.Logger)
	c.ShelfID = shelfID
	c.OwnerID = ownerID

	index := collection.NewIndex(scope, server.DBRegistry, server.Logger)
	index.ShelfID = shelfID
	index.OwnerID = ownerID

	err = index.Save(c, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	} else {
		response.Id = c.ID.String()
	}

	return response, nil
}

// SaveCollection updates a collection
func SaveCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveCollectionRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling save collection request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	shelfID, err := uuid.FromString(request.ShelfId)
	if err != nil {
		server.Logger.Debug("Invalid shelf id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	var ownerID uuid.UUID
	var scope collection.Scope
	if request.Scope == "account" {
		ownerID = server.Account.ID
		scope = collection.ScopeAccount
	} else if request.Scope == "user" {
		ownerID = server.Account.ActiveUser.ID
		scope = collection.ScopeUser
	} else {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	c := collection.New(t, scope, server.DBRegistry, server.Logger)
	c.ShelfID = shelfID
	c.OwnerID = ownerID
	c.Locked = request.Locked

	index := collection.NewIndex(scope, server.DBRegistry, server.Logger)
	index.ShelfID = shelfID
	index.OwnerID = ownerID

	err = index.Save(c, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}

// DeleteCollection deletes a collection
func DeleteCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.DeleteCollectionRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling delete collection request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	shelfID, err := uuid.FromString(request.ShelfId)
	if err != nil {
		server.Logger.Debug("Invalid shelf id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	var ownerID uuid.UUID
	var scope collection.Scope
	if request.Scope == "account" {
		ownerID = server.Account.ID
		scope = collection.ScopeAccount
	} else if request.Scope == "user" {
		ownerID = server.Account.ActiveUser.ID
		scope = collection.ScopeUser
	} else {
		return response, nil
	}

	c := collection.New(nil, scope, server.DBRegistry, server.Logger)
	c.ID = id
	c.ShelfID = shelfID
	c.OwnerID = ownerID

	index := collection.NewIndex(scope, server.DBRegistry, server.Logger)
	index.ShelfID = shelfID
	index.OwnerID = ownerID

	err = index.Delete(c, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}
