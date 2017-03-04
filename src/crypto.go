package main

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const (
	// KeySize is the number of bytes used for encryption keys
	KeySize = 32
	// NonceSize is the number of bytes used for the unique nonce
	NonceSize = 24
	// SaltSize is the size of the scrypt salt.
	SaltSize = 32
	// ScryptNVal is the N parameter value for scrypt
	ScryptNVal = 16384
	// ScryptRVal is the R parameter value for scrypt
	ScryptRVal = 8
	// ScryptPVal is the P parameter value for scrypt
	ScryptPVal = 1
	// ScryptKeySize is the size of the derived scrypt key
	ScryptKeySize = 32
)

var (
	// ErrEncrypt is the default error for encryption failures
	ErrEncrypt = errors.New("encryption failed")
	// ErrDecrypt is the default error for decryption failures
	ErrDecrypt = errors.New("decryption failed")
)

// RandBytes generates a sequence of random bytes
func RandBytes(size int) ([]byte, error) {
	r := make([]byte, size)
	_, err := rand.Read(r)
	return r, err
}

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

// DeriveKey takes a passphrase and derives a key from it
func DeriveKey(passphrase []byte, salt []byte) (*[KeySize]byte, error) {
	// N.B. - make sure to use subtle.ConstantTimeCompare when authenticating keys
	dk, err := scrypt.Key([]byte(passphrase), salt, ScryptNVal, ScryptRVal, ScryptPVal, ScryptKeySize)
	var key = new([KeySize]byte)
	copy(key[:], dk)
	Zero(dk)
	return key, err
}

// DeriveKeyAndSalt takes a passphrase and derives a key and salt from it
func DeriveKeyAndSalt(passphrase []byte) (*[KeySize]byte, []byte, error) {
	salt, err := RandBytes(SaltSize)
	if err != nil {
		return nil, nil, ErrEncrypt
	}
	key, err := DeriveKey(passphrase, salt)
	if err != nil {
		return nil, nil, ErrEncrypt
	}
	return key, salt, nil
}

// DeriveSaltedKey takes a passphrase and derives a key from it
// Additionally, the salt is embedded in the returned key
func DeriveSaltedKey(passphrase []byte) ([]byte, error) {
	key, salt, err := DeriveKeyAndSalt(passphrase)
	if err != nil {
		return nil, ErrEncrypt
	}
	out := key[0:len(key)]
	Zero(key[:])
	out = append(salt, out...)
	return out, nil
}

// ExtractSalt separates the salt and key into two pieces
func ExtractSalt(key []byte) ([]byte, *[KeySize]byte) {
	salt := key[:SaltSize]
	var onlyKey = new([KeySize]byte)
	copy(onlyKey[:], key[SaltSize:])
	return salt, onlyKey
}

// Seal secures a message using a passphrase.
func Seal(pass, message []byte) ([]byte, error) {
	salt, err := RandBytes(SaltSize)
	if err != nil {
		return nil, ErrEncrypt
	}

	key, err := DeriveKey(pass, salt)
	if err != nil {
		return nil, ErrEncrypt
	}
	out, err := Encrypt(key, message)
	Zero(key[:]) // Zero key immediately after
	if err != nil {
		return nil, ErrEncrypt
	}

	out = append(salt, out...)
	return out, nil
}

const overhead = SaltSize + secretbox.Overhead + NonceSize

// Open recovers a message encrypted using a passphrase.
func Open(pass, message []byte) ([]byte, error) {
	if len(message) < overhead {
		return nil, ErrDecrypt
	}

	key, err := DeriveKey(pass, message[:SaltSize])
	if err != nil {
		return nil, ErrDecrypt
	}
	out, err := Decrypt(key, message[SaltSize:])
	Zero(key[:]) // Zero key immediately after
	if err != nil {
		return nil, ErrDecrypt
	}

	return out, nil
}

// Zero attempts to zeroise its input.
func Zero(in []byte) {
	for i := 0; i < len(in); i++ {
		in[i] = 0
	}
}
