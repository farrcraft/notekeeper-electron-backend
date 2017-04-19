package rpc

import (
	"../notebook"
	messages "../proto"
	"github.com/golang/protobuf/proto"
)

// CreateNotebook is the RPC method to create a new notebook
func CreateNotebook(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: newResponseHeader(),
	}

	notebook := notebook.NewNotebook(rpc.DBFactory, rpc.Logger)

	err := notebook.Save(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		setInternalError(response.Header, err)
		return response, nil
	}

	response.Id = notebook.ID.String()

	return response, nil
}
