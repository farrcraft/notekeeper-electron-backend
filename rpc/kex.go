package rpc

import (
	"encoding/base64"

	"github.com/agl/ed25519"
)

// CreateSignature creates a signature for a response body
func (rpc *Server) CreateSignature(response []byte, context *RequestContext) string {
	signature := ed25519.Sign(context.Token.SignPrivateKey, response)

	sig := base64.StdEncoding.EncodeToString(signature[:])
	return sig

	/*
		rpc.SendCounter++
		response.Sequence = rpc.SendCounter
		return true
	*/
}

// VerifyRequest uses the client's public key to verify the message signature
func (rpc *Server) VerifyRequest(message []byte, sig []byte, context *RequestContext) bool {
	var signature [ed25519.SignatureSize]byte
	copy(signature[:], sig)
	ok := ed25519.Verify(context.Token.VerifyPublicKey, message, &signature)
	if !ok {
		rpc.Logger.Warn("Request payload signature could not be verified. key [", context.Token.VerifyPublicKey[:], "] message [", message, "] signature")
		return false
	}
	return true
}
