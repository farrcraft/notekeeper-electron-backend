package rpc

import (
	"../codes"
	"../notebook"
)

// CreateNotebook is the RPC method to create a new notebook
func CreateNotebook(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}
	notebook := notebook.NewNotebook(rpc.DB, rpc.Logger)
	err := notebook.Save(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
		return response, nil
	}
	return response, nil
}
