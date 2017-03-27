package rpc

import (
	"path/filepath"
	"time"

	"../codes"
	pb "../proto"
	"../uistate"
	"github.com/boltdb/bolt"
	"golang.org/x/net/context"
	//	"google.golang.org/grpc/credentials"
)

const (
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
)

// OpenMasterDb opens the master database in the requested directory
func (rpc *Server) OpenMasterDb(ctx context.Context, request *pb.OpenMasterDbRequest) (*pb.StatusResponse, error) {
	// need to close any existing db
	if rpc.DB != nil {
		rpc.DB.Close()
		rpc.DB = nil
	}
	rpc.DataPath = filepath.Clean(request.Path)
	// This is the master index db
	// There are additional databases where actual notebook data is stored (one DB file per account)
	fileName := filepath.Join(rpc.DataPath, MasterDbFile)
	rpc.Logger.Info("Opening master db file [", fileName, "]")
	var err error
	response := &pb.StatusResponse{}
	rpc.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		response.Status = codes.StatusError
		response.Code = int32(codes.ErrorMasterDbOpen)
		return response, nil
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err = state.Create()
	if err != nil {
		response.Status = codes.StatusError
		response.Code = int32(codes.ErrorCreateUIState)
		return response, nil
	}
	response.Status = codes.StatusOK
	return response, nil
}
