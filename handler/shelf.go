package handler

import (
	"notekeeper-electron-backend/codes"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"
	"notekeeper-electron-backend/shelf"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

func strToShelfScope(server *rpc.Server, s string, id uuid.UUID) (shelf.Scope, bool) {
	var scope shelf.Scope
	if s == "account" {
		if server.Account.ID != id {
			return scope, false
		}
		scope = shelf.ScopeAccount
	} else if s == "user" {
		if server.Account.ActiveUser.ID != id {
			return scope, false
		}
		scope = shelf.ScopeUser
	} else {
		return scope, false
	}
	return scope, true
}

func getShelves(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.GetShelvesResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetShelvesRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling get shelves request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Warn("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	shelfScope, ok := strToShelfScope(server, scope, id)
	if !ok {
		return response, nil
	}

	index := shelf.NewIndex(shelfScope, id, server.DBRegistry, server.Logger)
	err = index.LoadAll(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	for _, s := range index.Shelves {
		m := &messages.Shelf{
			Id:      s.ID.String(),
			Name:    rpc.TitleToMessage(s.Title),
			Default: s.Default,
			Locked:  s.Locked,
			Created: rpc.TimeToMessage(s.Created),
			Updated: rpc.TimeToMessage(s.Updated),
		}
		response.Shelves = append(response.Shelves, m)
	}

	return response, nil
}

func createShelf(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling create shelf request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Warn("Invalid owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	shelfScope, ok := strToShelfScope(server, scope, ownerID)
	if !ok {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	s, err := shelf.New(t, shelfScope, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating shelf - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	s.OwnerID = ownerID

	index := shelf.NewIndex(shelfScope, ownerID, server.DBRegistry, server.Logger)
	err = index.Save(s, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	} else {
		response.Id = s.ID.String()
	}

	return response, nil
}

// save an existing shelf
func saveShelf(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling save shelf request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	shelfScope, ok := strToShelfScope(server, scope, ownerID)
	if !ok {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	s, err := shelf.New(t, shelfScope, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating shelf - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	s.ID = id
	s.OwnerID = ownerID
	s.Locked = request.Locked

	index := shelf.NewIndex(shelfScope, ownerID, server.DBRegistry, server.Logger)
	err = index.Save(s, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}

func deleteShelf(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling delete shelf request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Warn("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	shelfScope, ok := strToShelfScope(server, scope, ownerID)
	if !ok {
		return response, nil
	}

	s, err := shelf.New(nil, shelfScope, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating shelf - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
	s.ID = id
	s.OwnerID = ownerID

	index := shelf.NewIndex(shelfScope, ownerID, server.DBRegistry, server.Logger)
	err = index.Delete(s)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}
