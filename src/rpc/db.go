package rpc

import (
	"../codes"
	"../db"
	messages "../proto"
	"../uistate"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// OpenMasterDb opens the master database in the requested directory
func OpenMasterDb(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: newResponseHeader(),
	}

	// need to close any existing db
	if rpc.DBFactory != nil {
		rpc.DBFactory.CloseAll()
	}

	request := messages.OpenMasterDbRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling open master db request - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	rpc.DBFactory = db.NewFactory(request.Path, rpc.Logger)

	// This is the master index db
	// There are additional databases where actual notebook data is stored
	db := rpc.DBFactory.DB(db.TypeMaster, uuid.Nil)
	rpc.Logger.Info("Opening master db file [", db.Filename, "]")
	err = db.Open()
	if err != nil {
		setRPCError(response.Header, codes.ErrorDbOpen)
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(db, rpc.Logger)
	err = state.Create()
	if err != nil {
		setInternalError(response.Header, err)
	}
	return response, nil
}
