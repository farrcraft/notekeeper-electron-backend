package rpc

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"crypto/rand"
	"crypto/tls"

	"strconv"

	"../account"
	"../db"
	"github.com/Sirupsen/logrus"
	"github.com/agl/ed25519"
	"github.com/golang/protobuf/proto"
)

// Handler is an RPC message handler
type Handler func(*Server, []byte) (proto.Message, error)

// Server is a RPC server instance
type Server struct {
	Logger          *logrus.Logger
	DBFactory       *db.Factory
	UserState       UserState
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
		UserState:   UserStateSignedOut,
	}
	return server
}

// RequestHeader contains the custom headers from a request
type RequestHeader struct {
	Signature []byte
	Method    string
	Sequence  int32
}

// VerifyHeaders checks that a request contains the correct headers &
// extracts their values into a working structure
func (rpc *Server) VerifyHeaders(req *http.Request) *RequestHeader {
	header := &RequestHeader{}

	header.Method = req.Header.Get("NoteKeeper-Request-Method")
	if header.Method == "" {
		rpc.Logger.Debug("Missing request method")
		return nil
	}

	// base64 encoded signature of the request body
	signature := req.Header.Get("NoteKeeper-Message-Signature")
	if signature == "" {
		rpc.Logger.Debug("Missing request signature")
		return nil
	}
	var err error
	header.Signature, err = base64.StdEncoding.DecodeString(signature)
	if err != nil {
		rpc.Logger.Debug("Error decoding request signature - ", err)
		return nil
	}

	seq := req.Header.Get("NoteKeeper-Message-Sequence")
	if seq == "" {
		rpc.Logger.Debug("Missing request sequence")
		return nil
	}
	parsedSeq, err := strconv.ParseInt(seq, 10, 32)
	if err != nil {
		rpc.Logger.Debug("Error decoding request sequence - ", err)
		return nil
	}
	header.Sequence = int32(parsedSeq)

	rpc.RecvCounter++
	if header.Method != "KeyExchange" && header.Sequence != rpc.RecvCounter {
		rpc.Logger.Debug("Invalid message sequence received. Expected [", rpc.RecvCounter, "] but got [", header.Sequence, "]")
		return nil
	}

	return header
}

// ServeHTTP handles HTTP requests
func (rpc *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// we only accept POST requests
	if req.Method != "POST" {
		rpc.Logger.Debug("Unexpected request method - ", req.Method)
		return
	}

	// we accept only one URL path of "/rpc"
	if req.URL.Path != "/rpc" {
		rpc.Logger.Debug("Unexpected request path - ", req.URL.Path)
		return
	}

	header := rpc.VerifyHeaders(req)
	if header == nil {
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rpc.Logger.Debug("Error reading request body - ", err)
		return
	}

	handler := rpc.FindHandler(header.Method)
	if handler == nil {
		rpc.Logger.Debug("Could not find handler for method - ", header.Method)
		return
	}

	// key exchange requests contain the key needed to do verification
	// so we need to defer until after the request has been handled
	if header.Method != "KeyExchange" {
		ok := rpc.VerifyRequest(body, header.Signature)
		if !ok {
			rpc.Logger.Debug("Message Verification failed")
			return
		}
	}

	handlerResponse, err := handler(rpc, body)
	if err != nil {
		return
	}

	if header.Method == "KeyExchange" {
		ok := rpc.VerifyRequest(body, header.Signature)
		if !ok {
			rpc.Logger.Debug("Message Verification failed")
			return
		}
	}

	responseData, err := proto.Marshal(handlerResponse)
	if err != nil {
		rpc.Logger.Debug("Error marshaling response - ", err)
		return
	}
	encodedData := base64.StdEncoding.EncodeToString(responseData)

	// set response headers
	responseSignature := rpc.CreateSignature(responseData)
	resp.Header().Set("NoteKeeper-Message-Signature", responseSignature)

	rpc.SendCounter++
	resp.Header().Set("NoteKeeper-Message-Sequence", strconv.FormatInt(int64(rpc.SendCounter), 10))
	// repackage request method header so client doesn't need to keep track of it
	resp.Header().Set("NoteKeeper-Request-Method", header.Method)

	// send response
	_, err = resp.Write([]byte(encodedData))
	if err != nil {
		rpc.Logger.Debug("Error writing response - ", err)
	}
}

// FindHandler matches a method name with a handler
func (rpc *Server) FindHandler(requestMethod string) Handler {
	for method, handler := range rpc.Handlers {
		if method == requestMethod {
			return handler
		}
	}
	return nil
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
	return true
}

// Stop performs shutdown routines before application termination
func (rpc *Server) Stop() {
	if rpc.DBFactory == nil {
		return
	}
	rpc.DBFactory.CloseAll()
	rpc.DBFactory = nil
}
