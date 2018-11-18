package main

import (
	"fmt"
	"os"

	"./handler"
	"./rpc"

	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	//	"google.golang.org/grpc/credentials"
	//	"golang.org/x/crypto/nacl/box"
)

const (
	// BackendPort is the GRPC service port
	BackendPort = "localhost:53017"
	// LogFile is the default log file name
	LogFile = "notekeeper.log"
)

// Backend is the main service type
type Backend struct {
	Logger   *logrus.Logger
	DB       *bbolt.DB // This is the master application DB
	RPC      *rpc.Server
	Status   chan string
	Shutdown chan bool
	//Account *Account
}

// NewBackend creates a new backend object
func NewBackend() *Backend {
	backend := &Backend{
		Logger:   logrus.New(),
		Status:   make(chan string),
		Shutdown: make(chan bool),
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

// Run is called when the application is started
func (backend *Backend) Run() {
	backend.RPC = rpc.NewServer(backend.Logger, backend.Status, backend.Shutdown)
	backend.RPC.RegisterHandlers(handler.Handlers())
	go backend.RPC.Start(BackendPort)
	for {
		select {
		case msg := <-backend.Status:
			fmt.Println(msg)
		case ok := <-backend.Shutdown:
			backend.Logger.Info("Shutting down service...")
			backend.RPC.Stop()
			if !ok {
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		}
	}
}
