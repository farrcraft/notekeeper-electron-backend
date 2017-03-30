package rpc

import (
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"

	"crypto/tls"

	"../account"
	"../appdir"
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

func (rpc *Server) createCertificate() bool {
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		rpc.Logger.Debug("Error creating serial number - ", err)
		return false
	}
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		rpc.Logger.Debug("Error generating private key - ", err)
		return false
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"NoteKeeper.io"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		IsCA:      true,
		KeyUsage:  x509.KeyUsageCertSign,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	keyPemBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	if err != nil {
		rpc.Logger.Debug("Error creating certificate - ", err)
		return false
	}
	certPemBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}

	certPath := filepath.Join(appdir.AppDataPath(), "certificate")
	certOut, err := os.Create(certPath)
	if err != nil {
		rpc.Logger.Debug("Error creating certificate file - ", err)
		return false
	}
	pem.Encode(certOut, &certPemBlock)
	certOut.Close()

	rpc.Certificate, err = tls.X509KeyPair(pem.EncodeToMemory(&certPemBlock), pem.EncodeToMemory(&keyPemBlock))
	if err != nil {
		rpc.Logger.Debug("Error converting certificate - ", err)
		return false
	}

	return true
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
	if err != nil {
		rpc.Logger.Debug("Credentials error - ", err)
		return false
	}
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
