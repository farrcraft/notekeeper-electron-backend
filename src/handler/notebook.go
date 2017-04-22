package handler

import (
	"../notebook"
	messages "../proto"
	"../rpc"
	"github.com/golang/protobuf/proto"
)

// CreateNotebook is the RPC method to create a new notebook
func CreateNotebook(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	notebook := notebook.NewNotebook(server.DBFactory, server.Logger)

	err := notebook.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	response.Id = notebook.ID.String()

	return response, nil
}
