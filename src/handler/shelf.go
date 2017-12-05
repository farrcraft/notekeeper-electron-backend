package handler

import (
	"../codes"
	messages "../proto"
	"../rpc"
	"../shelf"
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

// GetShelves gets a list of shelves
func GetShelves(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.GetShelvesResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetShelvesRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling get shelves request - ", err)
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

	scope, ok := strToShelfScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	index := shelf.NewIndex(scope, server.DBRegistry, server.Logger)
	err = index.LoadAll(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	for _, s := range index.Shelves {
		m := &messages.Shelf{
			Id:      s.ID.String(),
			Name:    rpc.TitleToMessage(s.Title),
			Scope:   request.Scope,
			Default: s.Default,
			Locked:  s.Locked,
			Created: rpc.TimeToMessage(s.Created),
			Updated: rpc.TimeToMessage(s.Updated),
		}
		response.Shelves = append(response.Shelves, m)
	}

	return response, nil
}

// CreateShelf saves a new shelf
func CreateShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create shelf request - ", err)
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

	scope, ok := strToShelfScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	s := shelf.New(t, scope, server.DBRegistry, server.Logger)

	index := shelf.NewIndex(scope, server.DBRegistry, server.Logger)
	err = index.Save(s, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	} else {
		response.Id = s.ID.String()
	}

	return response, nil
}

// SaveShelf saves an existing shelf
func SaveShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling save shelf request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := strToShelfScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	s := shelf.New(t, scope, server.DBRegistry, server.Logger)
	s.ID = id
	s.Locked = request.Locked

	index := shelf.NewIndex(scope, server.DBRegistry, server.Logger)
	err = index.Save(s, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}

// DeleteShelf deletes an existing shelf
func DeleteShelf(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling delete shelf request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := strToShelfScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	s := shelf.New(nil, scope, server.DBRegistry, server.Logger)
	s.ID = id

	index := shelf.NewIndex(scope, server.DBRegistry, server.Logger)
	err = index.Delete(s, server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}
