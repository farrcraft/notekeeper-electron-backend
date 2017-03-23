package rpc

import (
	"../codes"
	pb "../proto"
	"../uistate"
	"golang.org/x/net/context"
)

// UIState returns the the UI state as saved by the master DB
func (rpc *Server) UIState(ctx context.Context, request *pb.TokenRequest) (*pb.UIStateResponse, error) {
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err := state.Load()
	response := &pb.UIStateResponse{
		Status: &pb.StatusResponse{},
	}
	if err != nil {
		response.Status.Status = codes.StatusError
		response.Status.Code = int32(codes.ErrorLoadUIState)
		return response, nil
	}
	response.Status.Status = codes.StatusOK
	response.WindowWidth = state.WindowWidth
	response.WindowHeight = state.WindowHeight
	response.WindowXPosition = state.WindowXPosition
	response.WindowYPosition = state.WindowYPosition
	response.WindowMaximized = state.WindowMaximized
	response.WindowMinimized = state.WindowMinimized
	response.WindowFullscreen = state.WindowFullscreen
	response.DisplayWidth = state.DisplayWidth
	response.DisplayHeight = state.DisplayHeight
	response.DisplayXPosition = state.DisplayXPosition
	response.DisplayYPosition = state.DisplayYPosition

	return response, nil
}

// SaveUIState saves the current UI state to the master DB
func (rpc *Server) SaveUIState(ctx context.Context, request *pb.SaveUIStateRequest) (*pb.StatusResponse, error) {
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	state.WindowWidth = request.WindowWidth
	state.WindowHeight = request.WindowHeight
	state.WindowXPosition = request.WindowXPosition
	state.WindowYPosition = request.WindowYPosition
	state.WindowMaximized = request.WindowMaximized
	state.WindowMinimized = request.WindowMinimized
	state.WindowFullscreen = request.WindowFullscreen
	state.DisplayWidth = request.DisplayWidth
	state.DisplayHeight = request.DisplayHeight
	state.DisplayXPosition = request.DisplayXPosition
	state.DisplayYPosition = request.DisplayYPosition

	err := state.Save()
	response := &pb.StatusResponse{}
	if err != nil {
		response.Status = codes.StatusError
		response.Code = int32(codes.ErrorSaveUIState)
		return response, nil
	}
	response.Status = codes.StatusOK
	return response, nil
}
