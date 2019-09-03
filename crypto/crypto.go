package crypto

import (
	"crypto/rand"
	"io"

	"github.com/sirupsen/logrus"

	"notekeeper-electron-backend/codes"

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

// Context is a structure used for performing cryptography functions
type Context struct {
	Logger *logrus.Logger
}

// New creates a new Crypto object
func New(logger *logrus.Logger) *Context {
	c := &Context{
		Logger: logger,
	}
	return c
}

// RandBytes generates a sequence of random bytes
func (c *Context) RandBytes(size int) ([]byte, error) {
	r := make([]byte, size)
	_, err := rand.Read(r)
	if err != nil {
		c.Logger.Warn("Could not read random bytes - ", err)
		code := codes.New(codes.ScopeCrypto, codes.ErrorCrypto)
		return r, code
	}
	return r, nil
}

// GenerateKey creates a new random secret key.
func (c *Context) GenerateKey() (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		c.Logger.Warn("Could not read random key bytes - ", err)
		code := codes.New(codes.ScopeCrypto, codes.ErrorCrypto)
		return nil, code
	}

	return key, nil
}

// GenerateNonce creates a new random nonce.
func (c *Context) GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		c.Logger.Warn("Could not read random nonce bytes - ", err)
		code := codes.New(codes.ScopeCrypto, codes.ErrorCrypto)
		return nil, code
	}

	return nonce, nil
}

// Encrypt generates a random nonce and encrypts the input using
// NaCl's secretbox package. The nonce is prepended to the ciphertext.
// A sealed message will the same size as the original message plus
// secretbox.Overhead bytes long.
func (c *Context) Encrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	nonce, err := c.GenerateNonce()
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(nonce))
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, key)
	return out, nil
}

// Decrypt extracts the nonce from the ciphertext, and attempts to
// decrypt with NaCl's secretbox.
func (c *Context) Decrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	if len(message) < (NonceSize + secretbox.Overhead) {
		c.Logger.Warn("Message to short to decrypt")
		code := codes.New(codes.ScopeCrypto, codes.ErrorCrypto)
		return nil, code
	}

	var nonce [NonceSize]byte
	copy(nonce[:], message[:NonceSize])
	out, ok := secretbox.Open(nil, message[NonceSize:], &nonce, key)
	if !ok {
		code := codes.New(codes.ScopeCrypto, codes.ErrorDecrypt)
		return nil, code
	}

	return out, nil
}

// DeriveKey takes a passphrase and derives a key from it
func (c *Context) DeriveKey(passphrase []byte, salt []byte) (*[KeySize]byte, error) {
	// N.B. - make sure to use subtle.ConstantTimeCompare when authenticating keys
	dk, err := scrypt.Key([]byte(passphrase), salt, ScryptNVal, ScryptRVal, ScryptPVal, ScryptKeySize)
	var key = new([KeySize]byte)
	copy(key[:], dk)
	Zero(dk)
	if err != nil {
		c.Logger.Warn("Error deriving key - ", err)
		code := codes.New(codes.ScopeCrypto, codes.ErrorCrypto)
		return key, code
	}
	return key, nil
}

// EmbedSalt takes a key and salt and embeds the salt into the key
func (c *Context) EmbedSalt(key *[KeySize]byte, salt []byte) []byte {
	out := key[0:len(key)]
	out = append(salt, out...)
	return out
}

// DeriveKeyAndSalt takes a passphrase and derives a key and salt from it
func (c *Context) DeriveKeyAndSalt(passphrase []byte) (*[KeySize]byte, []byte, error) {
	salt, err := c.RandBytes(SaltSize)
	if err != nil {
		return nil, nil, err
	}
	key, err := c.DeriveKey(passphrase, salt)
	if err != nil {
		return nil, nil, err
	}
	return key, salt, nil
}

// DeriveSaltedKey takes a passphrase and derives a key from it
// Additionally, the salt is embedded in the returned key
func (c *Context) DeriveSaltedKey(passphrase []byte) ([]byte, error) {
	key, salt, err := c.DeriveKeyAndSalt(passphrase)
	if err != nil {
		return nil, err
	}
	out := key[0:len(key)]
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
func (c *Context) Seal(pass, message []byte) ([]byte, error) {
	salt, err := c.RandBytes(SaltSize)
	if err != nil {
		return nil, err
	}

	key, err := c.DeriveKey(pass, salt)
	if err != nil {
		return nil, err
	}
	out, err := c.Encrypt(key, message)
	Zero(key[:]) // Zero key immediately after
	if err != nil {
		return nil, err
	}

	out = append(salt, out...)
	return out, nil
}

const overhead = SaltSize + secretbox.Overhead + NonceSize

// Open recovers a message encrypted using a passphrase.
func (c *Context) Open(pass, message []byte) ([]byte, error) {
	if len(message) < overhead {
		c.Logger.Warn("Message too short to open")
		code := codes.New(codes.ScopeCrypto, codes.ErrorCrypto)
		return nil, code
	}

	key, err := c.DeriveKey(pass, message[:SaltSize])
	if err != nil {
		return nil, err
	}
	out, err := c.Decrypt(key, message[SaltSize:])
	Zero(key[:]) // Zero key immediately after
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Zero attempts to zeroise its input.
func Zero(in []byte) {
	for i := 0; i < len(in); i++ {
		in[i] = 0
	}
}
