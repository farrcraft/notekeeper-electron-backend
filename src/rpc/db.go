package rpc

import (
	"path/filepath"
	"time"

	"../codes"
	"../uistate"
	"github.com/boltdb/bolt"
	"github.com/mitchellh/mapstructure"
)

const (
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
)

type requestOpenMasterDb struct {
	Path string `mapstructure:"path"`
}

// OpenMasterDb opens the master database in the requested directory
func OpenMasterDb(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	// need to close any existing db
	if rpc.DB != nil {
		rpc.DB.Close()
		rpc.DB = nil
	}

	var request requestOpenMasterDb
	err := mapstructure.Decode(message.Payload, &request)
	if err != nil {
		rpc.Logger.Debug("Error decoding open master db request payload - ", err)
		response.Code = int(codes.ErrorMasterDbOpenDecode)
		response.Status = codes.StatusError
		return response, nil
	}

	rpc.DataPath = filepath.Clean(request.Path)
	// This is the master index db
	// There are additional databases where actual notebook data is stored (one DB file per account)
	fileName := filepath.Join(rpc.DataPath, MasterDbFile)
	rpc.Logger.Info("Opening master db file [", fileName, "]")
	rpc.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		response.Status = codes.StatusError
		response.Code = int(codes.ErrorMasterDbOpen)
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err = state.Create()
	if err != nil {
		response.Status = codes.StatusError
		response.Code = int(codes.ErrorCreateUIState)
		return response, nil
	}
	response.Status = codes.StatusOK
	return response, nil
}
