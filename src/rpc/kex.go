package rpc

import (
	"github.com/golang/protobuf/proto"

	"encoding/base64"

	"../codes"
	messages "../proto"
	"github.com/agl/ed25519"
)

// CreateSignature creates a signature for a response body
func (rpc *Server) CreateSignature(response []byte) string {
	signature := ed25519.Sign(rpc.SignPrivateKey, response)

	sig := base64.StdEncoding.EncodeToString(signature[:])
	return sig

	/*
		rpc.SendCounter++
		response.Sequence = rpc.SendCounter
		return true
	*/
}

// VerifyRequest uses the client's public key to verify the message signature
func (rpc *Server) VerifyRequest(message []byte, sig []byte) bool {
	var signature [ed25519.SignatureSize]byte
	copy(signature[:], sig)
	ok := ed25519.Verify(rpc.VerifyPublicKey, message, &signature)
	if !ok {
		rpc.Logger.Debug("Request payload signature could not be verified. key [", rpc.VerifyPublicKey[:], "] message [", message, "] signature")
		return false
	}
	return true
}

// KeyExchange performs a key exchange between client & server
func KeyExchange(rpc *Server, message []byte) (proto.Message, error) {
	response := &messages.KeyExchangeResponse{
		Header: newResponseHeader(),
	}

	request := messages.KeyExchangeRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling message - ", err)
		setRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// [FIXME] - assert request key length matches our target array size

	// client sent its own public key so we can verify requests it sends us later
	// this is a bit weak sauce wrt security since signature & verification key
	// are contained in the same message body, but it does give us assurance
	// that we at least have a functional verification key.
	rpc.VerifyPublicKey = new([ed25519.PublicKeySize]byte)
	copy(rpc.VerifyPublicKey[:], request.PublicKey)

	// send our own public key so client can verify our responses
	response.PublicKey = rpc.SignPublicKey[:]

	// reset sequence counters
	rpc.SendCounter = 0
	rpc.RecvCounter = 1

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
