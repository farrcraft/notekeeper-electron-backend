package handler

import (
	"notekeeper-electron-backend/codes"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"
	"notekeeper-electron-backend/uistate"
	"github.com/golang/protobuf/proto"
)

// OpenMasterDb opens the master database in the requested directory
func OpenMasterDb(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	// need to close any existing db
	if server.DBRegistry != nil {
		server.DBRegistry.CloseAll()
	}

	request := messages.OpenMasterDbRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling open master db request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	err = server.DBRegistry.OpenMaster(request.Path)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(server.DBRegistry, server.Logger)
	err = state.Create()
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}
	return response, nil
}
