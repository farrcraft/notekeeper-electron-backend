package rpc

import (
	"../codes"
	"../notebook"
	messages "../proto"
	"github.com/golang/protobuf/proto"
)

// CreateNotebook is the RPC method to create a new notebook
func CreateNotebook(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	notebook := notebook.NewNotebook(rpc.DB, rpc.Logger)

	err := notebook.Save(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
		return response, nil
	}

	response.Id = notebook.ID.String()

	return response, nil
}
