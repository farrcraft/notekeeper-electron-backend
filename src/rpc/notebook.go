package rpc

import (
	"golang.org/x/net/context"

	"../notebook"
	pb "../proto"
)

// CreateNotebook is the GRPC method to create a new notebook
func (rpc *Server) CreateNotebook(ctx context.Context, request *pb.CreateNotebookRequest) (*pb.CreateNotebookResponse, error) {
	notebook := notebook.NewNotebook(rpc.DB, rpc.Logger)
	err := notebook.Save(rpc.Account.ActiveUser.PassphraseKey)
	if err != nil {
		return nil, err
	}
	response := &pb.CreateNotebookResponse{
		Status: "OK",
		Id:     notebook.ID.String(),
	}
	return response, nil
}
