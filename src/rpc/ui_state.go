package rpc

import (
	pb "../proto"
	"../uistate"
	"golang.org/x/net/context"
)

// UIState returns the the UI state as saved by the master DB
func (rpc *Server) UIState(ctx context.Context, request *pb.TokenRequest) (*pb.UIStateResponse, error) {
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err := state.Load()
	if err != nil {
		return nil, err
	}
	response := &pb.UIStateResponse{
		WindowWidth:  state.WindowWidth,
		WindowHeight: state.WindowHeight,
	}
	return response, nil
}

// SaveUIState saves the current UI state to the master DB
func (rpc *Server) SaveUIState(ctx context.Context, request *pb.SaveUIStateRequest) (*pb.StatusResponse, error) {
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	state.WindowWidth = request.WindowWidth
	state.WindowHeight = request.WindowHeight
	err := state.Save()
	if err != nil {
		return nil, err
	}
	response := &pb.StatusResponse{}

	return response, nil
}
