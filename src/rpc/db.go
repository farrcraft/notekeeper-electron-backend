package rpc

import (
	"path/filepath"
	"time"

	"../codes"
	messages "../proto"
	"../uistate"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

const (
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
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
	if rpc.DB != nil {
		rpc.DB.Close()
		rpc.DB = nil
	}

	request := messages.OpenMasterDbRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling open master db request - ", err)
		response.Header.Code = int32(codes.ErrorMasterDbOpenDecode)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	rpc.DataPath = filepath.Clean(request.Path)
	// This is the master index db
	// There are additional databases where actual notebook data is stored (one DB file per account)
	fileName := filepath.Join(rpc.DataPath, MasterDbFile)
	rpc.Logger.Info("Opening master db file [", fileName, "]")
	rpc.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		response.Header.Code = int32(codes.ErrorMasterDbOpen)
		response.Header.Status = codes.StatusError
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err = state.Create()
	if err != nil {
		response.Header.Code = int32(codes.ErrorCreateUIState)
		response.Header.Status = codes.StatusError
		return response, nil
	}
	return response, nil
}
