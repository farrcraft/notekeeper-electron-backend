package handler

import (
	"notekeeper-electron-backend/codes"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"

	"github.com/agl/ed25519"
	"github.com/golang/protobuf/proto"
)

// KeyExchange performs a key exchange between client & server
func KeyExchange(server *rpc.Server, message []byte) (proto.Message, error) {
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

	// [FIXME] - assert request key length matches our target array size

	// client sent its own public key so we can verify requests it sends us later
	// this is a bit weak sauce wrt security since signature & verification key
	// are contained in the same message body, but it does give us assurance
	// that we at least have a functional verification key.
	server.VerifyPublicKey = new([ed25519.PublicKeySize]byte)
	copy(server.VerifyPublicKey[:], request.PublicKey)

	// send our own public key so client can verify our responses
	response.PublicKey = server.SignPublicKey[:]

	// reset sequence counters
	server.SendCounter = 0
	server.RecvCounter = 1

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
