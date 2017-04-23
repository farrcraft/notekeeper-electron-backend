package handler

import (
	"../codes"
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

	return response, nil
}
