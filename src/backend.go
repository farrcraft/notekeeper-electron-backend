package main

import (
	pb "./proto"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"time"
)

const (
	BACKEND_PORT   = ":53017"
	LOG_LEVEL      = "DEBUG"
	LOG_FILE       = "notekeeper.log"
	MASTER_DB_FILE = "notekeeper.db"
)

type Backend struct {
	Logger *logrus.Logger
	DB     *bolt.DB // This is the master application DB
}

func NewBackend() *Backend {
	backend := &Backend{
		Logger: logrus.New(),
	}

	backend.Logger.Formatter = &logrus.JSONFormatter{}
	level, err := logrus.ParseLevel(LOG_LEVEL)
	if err != nil {
		fmt.Println("Invalid log level [", LOG_LEVEL, "]")
		os.Exit(1)
	}
	backend.Logger.Level = level

	var file *os.File
	file, err = os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println("Unable to open log file [", LOG_FILE, "] - ", err)
		os.Exit(1)
	}
	backend.Logger.Out = file

	// This is the master index db
	// There are additional databases where actual notebook data is stored (one DB file per account)
	backend.DB, err = bolt.Open(MASTER_DB_FILE, 0600, &bolt.Options{Timeout: 1 * time.second})
	if err != nil {
		backend.Logger.Error("Unable to open DB - ", err)
		os.Exit(1)
	}

	return backend
}

func (backend *Backend) Shutdown() {
	backend.DB.Close()
}

func (backend *Backend) Run() {
	listener, err := net.Listen("tcp", BACKEND_PORT)
	if err != nil {
		backend.Logger.Error("Listen error - ", err)
		os.Exit(1)
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

	pb.RegisterBackendServer(server, backend)

	reflection.Register(s)

	err = server.Serve(listener)
	if err != nil {
		backend.Logger.Error("Server error - ", err)
		os.Exit(1)
	}
}

func (backend *Backend) CreateNotebook(ctx context.Context, request *pb.CreateNotebookRequest) (*pb.CreateNotebookResponse, error) {
	id := uuid.NewV4()
	response := &pb.CreateNotebookResponse{
		Status: "OK",
		Id:     id,
	}
	// [FIXME] - encrypt requested notebook name value
	//ciphertext, _ := cryptosecretbox.CryptoSecretBox([]byte(request.Name), nonce, key)

	// [FIXME] - need to use account-specific DB
	backend.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("notebooks"))
		if err != nil {
			backend.Logger.Error("Error creating notebook bucket - ", err)
			return err
		}
		// [FIXME] - use encrypted name value
		err = b.Put(id, []byte(request.Name))
		if err != nil {
			backend.Logger.Error("Error saving notebook - ", err)
			return err
		}
		return nil
	})
	return response, nil
}
