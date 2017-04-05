package rpc

import (
	"encoding/json"

	"../codes"
	"github.com/agl/ed25519"
	"github.com/mitchellh/mapstructure"
)

// SignResponse adds a signature to the response
func (rpc *Server) SignResponse(response *Response) bool {
	// we are signing the payload not the response envelope itself
	payload, err := json.Marshal(response.Payload)
	if err != nil {
		rpc.Logger.Debug("Error marshaling response payload for signing - ", err)
		return false
	}

	signature := ed25519.Sign(rpc.SignPrivateKey, payload)
	response.Signature = string(signature[:])

	rpc.SendCounter++
	response.Sequence = rpc.SendCounter
	return true
}

// payloadKeyExchange is the payload for both request & response
type payloadKeyExchange struct {
	PublicKey string `json:"public_key" mapstructure:"public_key"`
}

// KeyExchange performs a key exchange between client & server
func KeyExchange(rpc *Server, message *Message) (*Response, error) {
	response := &Response{
		Code:   int(codes.ErrorOK),
		Status: codes.StatusOK,
	}

	var request payloadKeyExchange
	err := mapstructure.Decode(message.Payload, &request)
	if err != nil {
		rpc.Logger.Debug("Error decoding key exchange request payload - ", err)
		response.Code = int(codes.ErrorKeyExchangeDecode)
		response.Status = codes.StatusError
		return response, nil
	}

	// client sent its own public key so we can verify requests it sends us later
	rpc.VerifyPublicKey = new([ed25519.PublicKeySize]byte)
	copy(rpc.VerifyPublicKey[:], request.PublicKey)

	// send our own public key so client can verify our responses
	response.Payload = &payloadKeyExchange{
		PublicKey: string(rpc.SignPublicKey[:]),
	}

	return response, nil
}
