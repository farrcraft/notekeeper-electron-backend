package rpc

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"crypto/tls"

	"strconv"

	"notekeeper-electron-backend/account"
	"notekeeper-electron-backend/db"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// Handler is an RPC message handler
type Handler func(*Server, []byte, *RequestContext) (proto.Message, error)

// RequestHeader contains the custom headers from a request
type RequestHeader struct {
	Signature []byte
	Method    string
	Sequence  int32
	Token     string
}

// RequestContext provides contextual information about a request
type RequestContext struct {
	Token  *ClientToken
	Header *RequestHeader
}

// Server is a RPC server instance
type Server struct {
	Logger      *logrus.Logger
	DBRegistry  *db.Registry
	UserState   UserState
	Account     *account.Account
	Certificate tls.Certificate
	Status      chan string
	Shutdown    chan bool
	Handlers    map[string]Handler
	Clients     map[string]*ClientToken
}

// NewServer creates a new RPCServer instance
func NewServer(logger *logrus.Logger, Status chan string, Shutdown chan bool) *Server {
	server := &Server{
		Logger:     logger,
		Handlers:   make(map[string]Handler, 0),
		UserState:  UserStateSignedOut,
		Status:     Status,
		Shutdown:   Shutdown,
		Clients:    make(map[string]*ClientToken),
		DBRegistry: db.NewRegistry(logger),
	}
	return server
}

// VerifyHeaders checks that a request contains the correct headers &
// extracts their values into a working structure
func (rpc *Server) VerifyHeaders(req *http.Request, context *RequestContext) bool {

	context.Header = &RequestHeader{}

	context.Header.Method = req.Header.Get("NoteKeeper-Request-Method")
	if context.Header.Method == "" {
		rpc.Logger.Warn("Missing request method")
		return false
	}

	if context.Header.Method == "SERVICE-READY" {
		return true
	}

	// Token creation is part of key exchange, so it doesn't exist here yet
	if context.Header.Method != "KeyExchange" {
		context.Header.Token = req.Header.Get("NoteKeeper-Client-Token")
		if context.Header.Token == "" {
			rpc.Logger.Warn("Missing request client token")
			return false
		}
		var ok bool
		context.Token, ok = rpc.Clients[context.Header.Token]
		if !ok {
			rpc.Logger.Warn("Invalid client token")
			return false
		}
	}

	// base64 encoded signature of the request body
	signature := req.Header.Get("NoteKeeper-Message-Signature")
	if signature == "" {
		rpc.Logger.Warn("Missing request signature")
		return false
	}
	var err error
	context.Header.Signature, err = base64.StdEncoding.DecodeString(signature)
	if err != nil {
		rpc.Logger.Warn("Error decoding request signature - ", err)
		return false
	}

	seq := req.Header.Get("NoteKeeper-Message-Sequence")
	if seq == "" {
		rpc.Logger.Warn("Missing request sequence")
		return false
	}
	parsedSeq, err := strconv.ParseInt(seq, 10, 32)
	if err != nil {
		rpc.Logger.Warn("Error decoding request sequence - ", err)
		return false
	}
	context.Header.Sequence = int32(parsedSeq)

	context.Token.RecvCounter++
	if context.Header.Method != "KeyExchange" && context.Header.Sequence != context.Token.RecvCounter {
		rpc.Logger.Warn("Invalid message sequence received. Expected [", context.Token.RecvCounter, "] but got [", context.Header.Sequence, "]")
		return false
	}

	return true
}

// ServeHTTP handles HTTP requests
func (rpc *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	rpc.Logger.Debug("PING")

	// we only accept POST requests
	if req.Method != "POST" {
		rpc.Logger.Warn("Unexpected request method - ", req.Method)
		return
	}

	// we accept only one URL path of "/rpc"
	if req.URL.Path != "/rpc" {
		rpc.Logger.Warn("Unexpected request path - ", req.URL.Path)
		return
	}

	context := &RequestContext{}

	if !rpc.VerifyHeaders(req, context) {
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rpc.Logger.Warn("Error reading request body - ", err)
		return
	}

	rpc.Logger.Debug(context.Header.Method)
	if context.Header.Method == "SERVICE-READY" {
		_, err = resp.Write([]byte("OK"))
		rpc.Logger.Debug("ready!")
		if err != nil {
			rpc.Logger.Warn("Error writing response - ", err)
		}
		return
	}

	handler := rpc.FindHandler(context.Header.Method)
	if handler == nil {
		rpc.Logger.Warn("Could not find handler for method - ", context.Header.Method)
		return
	}

	// key exchange requests contain the key needed to do verification
	// so we need to defer until after the request has been handled
	if context.Header.Method != "KeyExchange" {
		ok := rpc.VerifyRequest(body, context.Header.Signature, context)
		if !ok {
			rpc.Logger.Warn("Message Verification failed")
			return
		}
	}

	handlerResponse, err := handler(rpc, body, context)
	if err != nil {
		return
	}

	if context.Header.Method == "KeyExchange" {
		ok := rpc.VerifyRequest(body, context.Header.Signature, context)
		if !ok {
			rpc.Logger.Warn("Message Verification failed")
			return
		}
	}

	responseData, err := proto.Marshal(handlerResponse)
	if err != nil {
		rpc.Logger.Warn("Error marshaling response - ", err)
		return
	}
	encodedData := base64.StdEncoding.EncodeToString(responseData)

	// set response headers
	responseSignature := rpc.CreateSignature(responseData, context)
	resp.Header().Set("NoteKeeper-Message-Signature", responseSignature)

	context.Token.SendCounter++
	resp.Header().Set("NoteKeeper-Message-Sequence", strconv.FormatInt(int64(context.Token.SendCounter), 10))
	// repackage request method header so client doesn't need to keep track of it
	resp.Header().Set("NoteKeeper-Request-Method", context.Header.Method)

	// send response
	_, err = resp.Write([]byte(encodedData))
	if err != nil {
		rpc.Logger.Warn("Error writing response - ", err)
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

	conn, err := net.Listen("tcp", port)
	if err != nil {
		rpc.Logger.Warn("Listen error - ", err)
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

	// send a token to stdout so the frontend knows the backend is done initializing
	rpc.Status <- "NOTEKEEPER_SERVICE_READY"

	server.Serve(tlsListener)
	return true
}

// Stop performs shutdown routines before application termination
func (rpc *Server) Stop() {
	if rpc.DBRegistry == nil {
		return
	}
	rpc.DBRegistry.CloseAll()
	rpc.DBRegistry = nil
}
