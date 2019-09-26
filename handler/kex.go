package handler

import (
	"notekeeper-electron-backend/codes"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"

	"github.com/agl/ed25519"
	"github.com/golang/protobuf/proto"
)

// KeyExchange performs a key exchange between client & server
func KeyExchange(server *rpc.Server, message []byte, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.KeyExchangeResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.KeyExchangeRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling message - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new client token
	context.Token, err = rpc.NewClientToken(server.Logger)
	if err != nil {
		rpc.SetRPCError(response.Header, codes.ErrorCrypto)
		return response, nil
	}
	// [FIXME] - assert request key length matches our target array size

	// client sent its own public key so we can verify requests it sends us later
	// this is a bit weak sauce wrt security since signature & verification key
	// are contained in the same message body, but it does give us assurance
	// that we at least have a functional verification key.
	context.Token.VerifyPublicKey = new([ed25519.PublicKeySize]byte)
	copy(context.Token.VerifyPublicKey[:], request.PublicKey)

	// send our own public key so client can verify our responses
	response.PublicKey = context.Token.SignPublicKey[:]
	// the client will also need to keep track of its identifying token for future requests
	response.Token = context.Token.Token

	// reset sequence counters
	context.Token.SendCounter = 0
	context.Token.RecvCounter = 1

	server.Clients[context.Token.Token] = context.Token

	return response, nil
	/*
		ok := rpc.VerifyRequest(message)
		if !ok {
			response.Code = int(codes.ErrorVerifyRequestSignature)
			response.Status = codes.StatusError
			return response, nil
		}
	*/
}
