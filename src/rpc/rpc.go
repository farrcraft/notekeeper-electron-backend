package rpc

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"crypto/rand"
	"crypto/tls"

	"../account"
	"github.com/Sirupsen/logrus"
	"github.com/agl/ed25519"
	"github.com/boltdb/bolt"
)

// Message represents an RPC message
type Message struct {
	Method    string      `json:"method"`
	Signature string      `json:"signature"`
	Sequence  int32       `json:"sequence"`
	Payload   interface{} `json:"payload"`
}

// Response is an RPC response
type Response struct {
	Status    string      `json:"status"`
	Code      int         `json:"code"`
	Signature string      `json:"signature"`
	Sequence  int32       `json:"sequence"`
	Payload   interface{} `json:"payload"`
}

// Handler is an RPC message handler
type Handler func(*Server, *Message) (*Response, error)

// Server is a RPC server instance
type Server struct {
	Logger          *logrus.Logger
	DB              *bolt.DB // This is the master application DB
	DataPath        string
	Account         *account.Account
	Certificate     tls.Certificate
	Handlers        map[string]Handler
	RecvCounter     int32
	SendCounter     int32
	SignPublicKey   *[ed25519.PublicKeySize]byte
	SignPrivateKey  *[ed25519.PrivateKeySize]byte // Key used for signing responses
	VerifyPublicKey *[ed25519.PublicKeySize]byte  // Key used for verifying requests
}

// NewServer creates a new RPCServer instance
func NewServer(logger *logrus.Logger) *Server {
	server := &Server{
		Logger:      logger,
		Handlers:    make(map[string]Handler, 0),
		RecvCounter: 0,
		SendCounter: 0,
	}
	server.RegisterHandlers()
	return server
}

// RegisterHandlers registers all of the RPC handlers
func (rpc *Server) RegisterHandlers() {
	rpc.Handlers["KeyExchange"] = KeyExchange
	rpc.Handlers["MasterDb::open"] = OpenMasterDb

	rpc.Handlers["Account::create"] = CreateAccount
	rpc.Handlers["Account::unlock"] = UnlockAccount
	rpc.Handlers["Account::signin"] = SigninAccount
	rpc.Handlers["Account::signout"] = SignoutAccount
	rpc.Handlers["Account::lock"] = LockAccount

	rpc.Handlers["AccountState::get"] = GetAccountState

	rpc.Handlers["UIState::load"] = LoadUIState
	rpc.Handlers["UIState::save"] = SaveUIState
}

// ServeHTTP handles HTTP requests
func (rpc *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rpc.Logger.Debug("Unexpected request method - ", req.Method)
		return
	}

	// check req.URL to make sure it matches "/rpc"
	if req.URL.Path[1:] != "rpc" {
		rpc.Logger.Debug("Unexpected request path - ", req.URL.Path)
		return
	}

	message := &Message{}
	/*
		This is a short circuit for debugging raw input:
		body, _ := ioutil.ReadAll(req.Body)
		rpc.Logger.Debug(string(body))
	*/
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(message)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling request - ", err)
		return
	}

	// if this is a key exchange request, ignore the sequence
	// the internal sequence counters will be reset during the exchange process anyway
	if message.Method != "KeyExchange" {
		rpc.RecvCounter++
		if message.Sequence != rpc.RecvCounter {
			rpc.Logger.Debug("Invalid message sequence received. Expected [", rpc.RecvCounter, "] but got [", message.Sequence, "]")
			return
		}
	}

	foundHandler := false
	for method, handler := range rpc.Handlers {
		if method == message.Method {
			foundHandler = true
			handlerResponse, err := handler(rpc, message)
			if err != nil {
				return
			}

			ok := rpc.SignResponse(handlerResponse)
			if !ok {
				return
			}

			encoder := json.NewEncoder(resp)
			err = encoder.Encode(handlerResponse)
			if err != nil {
				rpc.Logger.Debug("Error marshaling handler response - ", err)
				return
			}
			break
		}
	}

	if !foundHandler {
		rpc.Logger.Debug("Could not find handler for method - ", message.Method)
		return
	}
}

// Start an RPC Server
func (rpc *Server) Start(port string) bool {
	ok := rpc.createCertificate()
	if !ok {
		return false
	}

	var err error
	rpc.SignPublicKey, rpc.SignPrivateKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		rpc.Logger.Debug("Error generating signing keys - ", err)
		return false
	}

	conn, err := net.Listen("tcp", port)
	if err != nil {
		rpc.Logger.Debug("Listen error - ", err)
		return false
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{rpc.Certificate},
	}
	tlsListener := tls.NewListener(conn, tlsConfig)
	writer := rpc.Logger.Writer()
	defer writer.Close()
	server := &http.Server{
		Addr:     port,
		Handler:  rpc,
		ErrorLog: log.New(writer, "", 0),
	}
	rpc.Logger.Debug("RPC listening on port [", port, "]")
	server.Serve(tlsListener)
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
	*/
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
