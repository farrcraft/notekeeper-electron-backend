package rpc

import (
	"golang.org/x/net/context"

	pb "../proto"
)

// CreateNotebook is the GRPC method to create a new notebook
func (rpc *RPCServer) CreateNotebook(ctx context.Context, request *pb.CreateNotebookRequest) (*pb.CreateNotebookResponse, error) {
	notebook := NewNotebook(backend.Account)
	err := notebook.Save()
	if err != nil {
		return nil, err
	}
	response := &pb.CreateNotebookResponse{
		Status: "OK",
		Id:     notebook.ID.String(),
	}
	return response, nil
}
