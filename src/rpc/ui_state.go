package rpc

import (
	"../codes"
	"../uistate"
	"github.com/mitchellh/mapstructure"
)

// LoadUIState returns the the UI state as saved by the master DB
func LoadUIState(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err := state.Load()
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
		return response, nil
	}
	response.Payload = state

	return response, nil
}

// SaveUIState saves the current UI state to the master DB
func SaveUIState(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err := mapstructure.Decode(message.Payload, state)
	if err != nil {
		rpc.Logger.Debug("Error decoding save ui state request - ", err)
		response.Status = codes.StatusError
		response.Code = int(codes.ErrorSaveUIStateDecode)
		return response, nil
	}

	err = state.Save()
	if err != nil {
		code := codes.ToInternalError(err)
		response.Status = code.Error()
		response.Code = int(code.Code)
		return response, nil
	}

	return response, nil
}
