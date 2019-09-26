package rpc

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/agl/ed25519"
	"github.com/sirupsen/logrus"
)

// ClientToken identifies a client that can communicate with the server
// Each client has its own set of counters and signing keys
type ClientToken struct {
	Token           string
	RecvCounter     int32
	SendCounter     int32
	SignPublicKey   *[ed25519.PublicKeySize]byte
	SignPrivateKey  *[ed25519.PrivateKeySize]byte // Key used for signing responses
	VerifyPublicKey *[ed25519.PublicKeySize]byte  // Key used for verifying requests
}

// NewClientToken creates a new ClientToken
func NewClientToken(logger *logrus.Logger) (*ClientToken, error) {
	client := &ClientToken{
		RecvCounter: 0,
		SendCounter: 0,
	}

	// The identifier token is just a url encoded random string
	tokenLength := 32
	rb := make([]byte, tokenLength)
	_, err := rand.Read(rb)
	if err != nil {
		logger.Warn("Error creating client token - ", err)
		return client, err
	}
	client.Token = base64.URLEncoding.EncodeToString(rb)

	client.SignPublicKey, client.SignPrivateKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		logger.Warn("Error generating signing keys - ", err)
		return client, err
	}

	return client, nil
}
