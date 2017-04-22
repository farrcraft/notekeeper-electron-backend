package handler

import (
	messages "../proto"
	"../rpc"

	"github.com/golang/protobuf/proto"
)

// GetCollections gets a list of collections
func GetCollections(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}
	return response, nil
}

// CreateCollection saves a new collection
func CreateCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}
	return response, nil
}

// SaveCollection updates a collection
func SaveCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}
	return response, nil
}

// DeleteCollection deletes a collection
func DeleteCollection(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}
	return response, nil
}
