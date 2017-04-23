package handler

import (
	"../codes"
	"../db"
	messages "../proto"
	"../rpc"
	"../uistate"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// OpenMasterDb opens the master database in the requested directory
func OpenMasterDb(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	// need to close any existing db
	if server.DBFactory != nil {
		server.DBFactory.CloseAll()
	}

	request := messages.OpenMasterDbRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling open master db request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	server.DBFactory = db.NewFactory(request.Path, server.Logger)

	// This is the master index db
	// There are additional databases where actual notebook data is stored
	db, err := server.DBFactory.DB(db.TypeMaster, uuid.Nil)
	server.Logger.Info("Opening master db file [", db.Filename, "]")
	if err != nil {
		rpc.SetRPCError(response.Header, codes.ErrorDbOpen)
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(db, server.Logger)
	err = state.Create()
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}
	return response, nil
}
