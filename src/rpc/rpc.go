package rpc

import (
	"net"

	"crypto/tls"

	"../account"
	pb "../proto"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

// Server is a gRPC server instance
type Server struct {
	Logger      *logrus.Logger
	DB          *bolt.DB // This is the master application DB
	DataPath    string
	Account     *account.Account
	Certificate tls.Certificate
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
		rpc.Logger.Debug("Listen error - ", err)
		return false
	}

	// make sure gRPC logging goes to our logger and not the default stderr
	grpclog.SetLogger(rpc.Logger)

	// [FIXME] - need to make this use TLS
	// can we just generate a certificate on the fly to use here?
	var opts []grpc.ServerOption
	ok := rpc.createCertificate()
	if !ok {
		return false
	}
	creds := credentials.NewServerTLSFromCert(&rpc.Certificate)
	/*
		if err != nil {
			rpc.Logger.Debug("Credentials error - ", err)
			return false
		}
	*/
	/*
		tlsConfig := &tls.Config{
			Rand:         rand.Reader,
			Certificates: []tls.Certificate{rpc.Certificate},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			},
		}
		creds := credentials.NewTLS(tlsConfig)
	*/
	/*
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			backend.Logger.Error("Credentials error - ", err)
			os.Exit(1)
		}
	*/
	opts = []grpc.ServerOption{grpc.Creds(creds)}
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
		rpc.DB = nil
	}
	if rpc.Account != nil && rpc.Account.DB != nil {
		rpc.Account.DB.Close()
		rpc.Account.DB = nil
	}
}
