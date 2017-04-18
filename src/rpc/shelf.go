package rpc

import (
	"../codes"
	messages "../proto"
	"../shelf"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// GetShelves gets a list of shelves
func GetShelves(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.GetShelvesResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	request := messages.GetShelvesRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling get shelves request - ", err)
		response.Header.Code = int32(codes.ErrorDecode)
		response.Header.Scope = int32(codes.ScopeRPC)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		rpc.Logger.Debug("Invalid id - ", err)
		response.Header.Code = int32(codes.ErrorDecode)
		response.Header.Scope = int32(codes.ScopeRPC)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	var scope shelf.Scope
	if request.Scope == "account" {
		if rpc.Account.ID != id {
			return response, nil
		}
		scope = shelf.ScopeAccount
	} else if request.Scope == "user" {
		if rpc.Account.ActiveUser.ID != id {
			return response, nil
		}
		scope = shelf.ScopeUser
	} else {
		rpc.Logger.Debug("Unrecognized scope")
		return response, nil
	}

	// create a new shelf instance to act as a proxy
	s := shelf.New(nil, scope, rpc.DBFactory, rpc.Logger)
	shelves, err := s.LoadAll(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Scope = int32(code.Scope)
		response.Header.Status = code.Error()
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

}

// SaveShelf saves an existing shelf
func SaveShelf(rpc *Server, message []byte) (proto.Message, error) {

}

// DeleteShelf deletes an existing shelf
func DeleteShelf(rpc *Server, message []byte) (proto.Message, error) {

}
