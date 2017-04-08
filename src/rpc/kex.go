package rpc

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"../codes"
	messages "../proto"
	"github.com/agl/ed25519"
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
	response.Signature = signature[:]

	rpc.SendCounter++
	response.Sequence = rpc.SendCounter
	return true
}

// VerifyRequest uses the client's public key to verify the message signature
func (rpc *Server) VerifyRequest(message []byte, sig []byte) bool {
	/*
		payload, err := json.Marshal(message.Payload)
		if err != nil {
			rpc.Logger.Debug("Error marshaling request payload for signature verification - ", err)
			return false
		}
	*/
	/*
		data, err := base64.StdEncoding.DecodeString(message.Signature)
		if err != nil {
			rpc.Logger.Debug("Error decoding signature - ", err)
			return false
		}
	*/
	var signature [ed25519.SignatureSize]byte
	/*
		var idx int
		// [FIXME] - assert message signature length matches our target array size
		for _, c := range message.Signature {
			signature[idx] = c
			idx++
		}
	*/
	copy(signature[:], sig)
	ok := ed25519.Verify(rpc.VerifyPublicKey, message, &signature)
	if !ok {
		rpc.Logger.Debug("Request payload signature could not be verified. key [", rpc.VerifyPublicKey[:], "] message [", message, "] signature")
		rpc.Logger.Debug(signature[:])
		//rpc.Logger.Debug(message.RawPayload)
		return false
	} else {
		rpc.Logger.Debug("Request payload signature was verified!")
	}
	return true
}

// payloadKeyExchange is the payload for both request & response
type payloadKeyExchange struct {
	PublicKey map[string]byte `json:"public_key" mapstructure:"public_key"`
}

// KeyExchange performs a key exchange between client & server
func KeyExchange(rpc *Server, message []byte) ([]byte, error) {
	response := messages.KeyExchangeResponse{
		Header: &messages.ResponseHeader{
			Code:   int32(codes.ErrorOK),
			Status: codes.StatusOK,
		},
	}

	request := messages.KeyExchangeRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		rpc.Logger.Debug("Error unmarshaling message - ", err)
		response.Header.Code = int32(codes.ErrorKeyExchangeDecode)
		response.Header.Status = codes.StatusError
		data, err := proto.Marshal(&response)
		if err != nil {
			rpc.Logger.Debug("Error marshaling error response")
		}
		return data, nil
	}

	/*
		response := &Response{
			Code:   int(codes.ErrorOK),
			Status: codes.StatusOK,
		}
	*/
	/*
		var request payloadKeyExchange

		data, err := base64.StdEncoding.DecodeString(message.RawPayload)
		if err != nil {
			rpc.Logger.Debug("Error decoding request payload - ", err)
			response.Header.Code = int32(codes.ErrorKeyExchangeDecode)
			response.Header.Status = codes.StatusError
			return response, nil
		}
		var payload map[string]interface{}
		err = json.Unmarshal(data, &payload)
		if err != nil {
			rpc.Logger.Debug("Error unmarshaling request payload - ", err)
			response.Code = int(codes.ErrorKeyExchangeDecode)
			response.Status = codes.StatusError
			return response, nil
		}
		err = mapstructure.Decode(payload, &request)
		if err != nil {
			rpc.Logger.Debug("Error decoding key exchange request payload - ", err)
			rpc.Logger.Debug(string(data))
			response.Code = int(codes.ErrorKeyExchangeDecode)
			response.Status = codes.StatusError
			return response, nil
		}
	*/
	// client sent its own public key so we can verify requests it sends us later
	rpc.VerifyPublicKey = new([ed25519.PublicKeySize]byte)
	copy(rpc.VerifyPublicKey[:], request.PublicKey)
	/*
		var idx int
		// [FIXME] - assert request key length matches our target array size
		for _, c := range request.PublicKey {
			rpc.VerifyPublicKey[idx] = c
			idx++
		}
	*/
	data, err := proto.Marshal(&response)
	if err != nil {
		rpc.Logger.Debug("Error marshaling error response")
	}
	return data, nil

	// this is a bit weak sauce wrt security since signature & verification key
	// are contained in the same message body, but it does give us assurance
	// that we at least have a functional verification key.
	/*
		ok := rpc.VerifyRequest(message)
		if !ok {
			response.Code = int(codes.ErrorVerifyRequestSignature)
			response.Status = codes.StatusError
			return response, nil
		}

		// reset sequence counters
		rpc.SendCounter = 0
		rpc.RecvCounter = 1

		// send our own public key so client can verify our responses
		response.Payload = &payloadKeyExchange{
		//PublicKey: rpc.SignPublicKey[:],
		}

		return response, nil
	*/
}
