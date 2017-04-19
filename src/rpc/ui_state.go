package rpc

import (
	"../codes"
	"../db"
	messages "../proto"
	"../uistate"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// LoadUIState returns the the UI state as saved by the master DB
func LoadUIState(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.LoadUIStateResponse{
		Header: newResponseHeader(),
	}

	db := rpc.DBFactory.Find(db.TypeMaster, uuid.Nil)
	state := uistate.NewUIState(db, rpc.Logger)
	err := state.Load()
	if err != nil {
		setInternalError(response.Header, err)
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
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	db := rpc.DBFactory.Find(db.TypeMaster, uuid.Nil)
	state := uistate.NewUIState(db, rpc.Logger)

	request := messages.SaveUIStateRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling save ui state request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
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
		setInternalError(response.Header, err)
		return response, nil
	}

	return response, nil
}
