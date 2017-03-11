package rpc

import (
	"fmt"
	"net"
	"path/filepath"
	"time"

	"../account"
	pb "../proto"
	"../uistate"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	//	"google.golang.org/grpc/credentials"
)

const (
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
)

// Server is a gRPC server instance
type Server struct {
	Logger  *logrus.Logger
	DB      *bolt.DB // This is the master application DB
	Account *account.Account
}

// NewServer creates a new RPCServer instance
func NewServer(logger *logrus.Logger) *Server {
	server := &Server{
		Logger: logger,
	}
	return server
}

// Start starts an RPCServer
func (rpc *Server) Start(port string) bool {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		rpc.Logger.Error("Listen error - ", err)
		return false
	}
	// [FIXME] - need to make this use TLS
	// can we just generate a certificate on the fly to use here?
	var opts []grpc.ServerOption
	/*
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			backend.Logger.Error("Credentials error - ", err)
			os.Exit(1)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	*/
	server := grpc.NewServer(opts...)

	pb.RegisterBackendServer(server, rpc)

	reflection.Register(server)

	rpc.Logger.Debug("RPC listening on port [", port, "]")
	err = server.Serve(listener)
	if err != nil {
		rpc.Logger.Error("Server error - ", err)
		return false
	}

	return true
}

// Stop performs shutdown routines before application termination
func (rpc *Server) Stop() {
	if rpc.DB != nil {
		rpc.DB.Close()
	}
	if rpc.Account != nil && rpc.Account.DB != nil {
		rpc.Account.DB.Close()
	}
}

// OpenMasterDb opens the master database in the requested directory
func (rpc *Server) OpenMasterDb(ctx context.Context, request *pb.OpenMasterDbRequest) (*pb.OpenMasterDbResponse, error) {
	// This is the master index db
	// There are additional databases where actual notebook data is stored (one DB file per account)
	fileName := fmt.Sprint(filepath.Clean(request.Path), filepath.Separator, MasterDbFile)
	rpc.Logger.Info("Opening master db file [", fileName, "]")
	var err error
	rpc.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		rpc.Logger.Error("Unable to open DB - ", err)
		return nil, err
	}

	// make sure DB has a default UIState saved
	state := uistate.NewUIState(rpc.DB, rpc.Logger)
	err = state.Create()
	if err != nil {
		return nil, err
	}

	response := &pb.OpenMasterDbResponse{}
	return response, nil
}
