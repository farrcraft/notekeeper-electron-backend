package rpc

import (
	"../codes"
	messages "../proto"
	"../shelf"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

func (rpc *Server) strToShelfScope(s string, id uuid.UUID) (shelf.Scope, bool) {
	var scope shelf.Scope
	if s == "account" {
		if rpc.Account.ID != id {
			return scope, false
		}
		scope = shelf.ScopeAccount
	} else if s == "user" {
		if rpc.Account.ActiveUser.ID != id {
			return scope, false
		}
		scope = shelf.ScopeUser
	} else {
		return scope, false
	}
	return scope, true
}

// GetShelves gets a list of shelves
func GetShelves(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.GetShelvesResponse{
		Header: newResponseHeader(),
	}

	request := messages.GetShelvesRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling get shelves request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		rpc.Logger.Debug("Invalid id - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := rpc.strToShelfScope(request.Scope, id)
	if !ok {
		return response, nil
	}

	// create a new shelf instance to act as a proxy
	s := shelf.New(nil, scope, rpc.DBFactory, rpc.Logger)
	shelves, err := s.LoadAll(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		setInternalError(response.Header, err)
		return response, nil
	}

	for _, s := range shelves {
		m := &messages.Shelf{
			Id:      s.ID.String(),
			Name:    titleToMessage(s.Title),
			Scope:   request.Scope,
			Default: s.Default,
			Locked:  s.Locked,
			Created: timeToMessage(s.Created),
			Updated: timeToMessage(s.Updated),
		}
		response.Shelves = append(response.Shelves, m)
	}

	return response, nil
}

// CreateShelf saves a new shelf
func CreateShelf(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: newResponseHeader(),
	}

	request := messages.CreateShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling create shelf request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		rpc.Logger.Debug("Invalid id - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := rpc.strToShelfScope(request.Scope, id)
	if !ok {
		return response, nil
	}

	t := messageToTitle(request.Name)
	s := shelf.New(t, scope, rpc.DBFactory, rpc.Logger)
	err = s.Save(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		setInternalError(response.Header, err)
	} else {
		response.Id = s.ID.String()
	}

	return response, nil
}

// SaveShelf saves an existing shelf
func SaveShelf(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	request := messages.SaveShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling save shelf request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	id, err := uuid.FromString(request.OwnerId)
	if err != nil {
		rpc.Logger.Debug("Invalid id - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := rpc.strToShelfScope(request.Scope, id)
	if !ok {
		return response, nil
	}

	t := messageToTitle(request.Name)
	s := shelf.New(t, scope, rpc.DBFactory, rpc.Logger)
	s.ID = id
	s.Locked = request.Locked
	err = s.Save(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		setInternalError(response.Header, err)
	}

	return response, nil
}

// DeleteShelf deletes an existing shelf
func DeleteShelf(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	request := messages.SaveShelfRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling delete shelf request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	id, err := uuid.FromString(request.OwnerId)
	if err != nil {
		rpc.Logger.Debug("Invalid id - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := rpc.strToShelfScope(request.Scope, id)
	if !ok {
		return response, nil
	}

	s := shelf.New(nil, scope, rpc.DBFactory, rpc.Logger)
	s.ID = id
	err = s.Delete()
	if err != nil {
		setInternalError(response.Header, err)
	}

	return response, nil
}
