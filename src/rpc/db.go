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
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	// need to close any existing db
	if rpc.DBFactory != nil {
		rpc.DBFactory.CloseAll()
	}

	request := messages.OpenMasterDbRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling open master db request - ", err)
		response.Header.Code = int32(codes.ErrorMasterDbOpenDecode)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	rpc.DBFactory = db.NewFactory(request.Path, rpc.Logger)

	// This is the master index db
	// There are additional databases where actual notebook data is stored
	db := rpc.DBFactory.DB(db.TypeMaster, uuid.Nil)
	rpc.Logger.Info("Opening master db file [", db.Filename, "]")
	err = db.Open()
	if err != nil {
		response.Header.Code = int32(codes.ErrorMasterDbOpen)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(db, rpc.Logger)
	err = state.Create()
	if err != nil {
		response.Header.Code = int32(codes.ErrorCreateUIState)
		response.Header.Status = codes.StatusError
		return response, nil
	}
	return response, nil
}
