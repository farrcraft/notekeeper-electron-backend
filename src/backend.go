package main

import (
	"fmt"
	"net"
	"os"

	pb "./proto"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	//	"google.golang.org/grpc/credentials"
	//	"golang.org/x/crypto/nacl/box"
	"github.com/kardianos/service"
)

const (
	// BackendPort is the GRPC service port
	BackendPort = ":53017"
	// LogLevel is the default log level
	LogLevel = "DEBUG"
	// LogFile is the default log file name
	LogFile = "notekeeper.log"
	// MasterDbFile is the core bolt database filename
	MasterDbFile = "notekeeper.db"
)

// Backend is the main service type
type Backend struct {
	Logger  *logrus.Logger
	DB      *bolt.DB // This is the master application DB
	Account *Account
}

// NewBackend creates a new backend object
func NewBackend() *Backend {
	backend := &Backend{
		Logger: logrus.New(),
	}

	backend.Logger.Formatter = &logrus.JSONFormatter{}
	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		fmt.Println("Invalid log level [", LogLevel, "]")
		os.Exit(1)
	}
	backend.Logger.Level = level

	var file *os.File
	file, err = os.OpenFile(LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println("Unable to open log file [", LogFile, "] - ", err)
		os.Exit(1)
	}
	backend.Logger.Out = file

	return backend
}

// Stop stops the backend service
func (backend *Backend) Stop(s service.Service) error {
	backend.Logger.Debug("Stopping service...")
	// Stop should not block. Return with a few seconds.
	backend.Shutdown()
	return nil
}

// Shutdown is called when the application is terminated
// Caveat - running via CLI on Windows under MSYS2 (e.g. babun) doesn't seem to capture ctrl^c
// So shutdown won't get called. Use normal CMD prompt or powershell in that scenario instead.
func (backend *Backend) Shutdown() {
	backend.Logger.Debug("Shutting down service...")
	if backend.DB != nil {
		backend.DB.Close()
	}
	if backend.Account != nil && backend.Account.DB != nil {
		backend.Account.DB.Close()
	}
}

// Start starts the backend service
func (backend *Backend) Start(svc service.Service) error {
	// Start should not block. Do the actual work async.
	go backend.Run()
	return nil
}

// Run is called when the application is started
func (backend *Backend) Run() {
	listener, err := net.Listen("tcp", BackendPort)
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

	reflection.Register(server)

	backend.Logger.Debug("Backend listening on port [", BackendPort, "]")
	err = server.Serve(listener)
	if err != nil {
		backend.Logger.Error("Server error - ", err)
		os.Exit(1)
	}
}
