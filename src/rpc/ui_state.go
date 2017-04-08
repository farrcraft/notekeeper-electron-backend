package rpc

import (
	"../codes"
	messages "../proto"
	"../uistate"
	"github.com/golang/protobuf/proto"
)

// LoadUIState returns the the UI state as saved by the master DB
func LoadUIState(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.LoadUIStateResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err := state.Load()
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
		return response, nil
	}

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
func SaveUIState(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.LoadUIStateResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	state := uistate.NewUIState(rpc.DB, rpc.Logger)

	request := messages.SaveUIStateRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling save ui state request - ", err)
		response.Header.Code = int32(codes.ErrorSaveUIStateDecode)
		response.Header.Status = codes.StatusError
		return response, nil
	}

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

	err = state.Save()
	if err != nil {
		code := codes.ToInternalError(err)
		response.Header.Code = int32(code.Code)
		response.Header.Status = code.Error()
		return response, nil
	}

	return response, nil
}
