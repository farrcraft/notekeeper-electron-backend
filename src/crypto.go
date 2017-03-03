package main

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// KeySize is the number of bytes used for encryption keys
	KeySize = 32
	// NonceSize is the number of bytes used for the unique nonce
	NonceSize = 24
)

var (
	// ErrEncrypt is the default error for encryption failures
	ErrEncrypt = errors.New("secret: encryption failed")
	// ErrDecrypt is the default error for decryption failures
	ErrDecrypt = errors.New("secret: decryption failed")
)

// GenerateKey creates a new random secret key.
func GenerateKey() (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateNonce creates a new random nonce.
func GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// Encrypt generates a random nonce and encrypts the input using
// NaCl's secretbox package. The nonce is prepended to the ciphertext.
// A sealed message will the same size as the original message plus
// secretbox.Overhead bytes long.
func Encrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, ErrEncrypt
	}

	out := make([]byte, len(nonce))
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, key)
	return out, nil
}

// Decrypt extracts the nonce from the ciphertext, and attempts to
// decrypt with NaCl's secretbox.
func Decrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	if len(message) < (NonceSize + secretbox.Overhead) {
		return nil, ErrDecrypt
	}

	var nonce [NonceSize]byte
	copy(nonce[:], message[:NonceSize])
	out, ok := secretbox.Open(nil, message[NonceSize:], &nonce, key)
	if !ok {
		return nil, ErrDecrypt
	}

	return out, nil
}
