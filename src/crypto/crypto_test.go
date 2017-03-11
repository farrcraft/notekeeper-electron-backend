package crypto

import (
	"bytes"
	"testing"
)

func TestCrypto(t *testing.T) {
	// RandBytes
	randomBytes, err := RandBytes(27)
	if err != nil {
		t.Error("Error creating random bytes - ", err)
	}
	if len(randomBytes) != 27 {
		t.Error("Unexpected random bytes length - ", len(randomBytes))
	}

	moreBytes, err := RandBytes(27)
	if err != nil {
		t.Error("Error creating more random bytes - ", err)
	}

	if bytes.Equal(randomBytes, moreBytes) {
		t.Error("Expected random bytes not to match")
	}

	// GenerateKey
	key, err := GenerateKey()
	if err != nil {
		t.Error("Expected to generate key - ", err)
	}
	if len(key) != KeySize {
		t.Error("Expected key length to be key size")
	}
	key2, err := GenerateKey()
	if err != nil {
		t.Error("Expected to generate key2 - ", err)
	}
	if bytes.Equal(key[:], key2[:]) {
		t.Error("Expected keys not to match")
	}

	// GenerateNonce
	nonce, err := GenerateNonce()
	if err != nil {
		t.Error("Expected to generate nonce - ", err)
	}
	if len(nonce) != NonceSize {
		t.Error("Expected nonce length to be nonce size")
	}
	nonce2, err := GenerateNonce()
	if err != nil {
		t.Error("Expected to generate nonce2 - ", err)
	}
	if bytes.Equal(nonce[:], nonce2[:]) {
		t.Error("Expected nonces not to match")
	}

	// Encrypt
	testMessage := []byte("This a test message.")
	encryptedMessage, err := Encrypt(key, testMessage)
	if err != nil {
		t.Error("Expected to encrypt message")
	}
	if bytes.Equal(testMessage, encryptedMessage) {
		t.Error("Expected encrypted message not to match original message")
	}

	encryptedMessage2, err := Encrypt(key2, testMessage)
	if err != nil {
		t.Error("Expected to encrypt message2")
	}
	if bytes.Equal(encryptedMessage, encryptedMessage2) {
		t.Error("Expected encrypted messages not to match")
	}

	// Decrypt
	decryptedMessage, err := Decrypt(key, encryptedMessage)
	if err != nil {
		t.Error("Expected to decrypt message - ", err)
	}
	if !bytes.Equal(testMessage, decryptedMessage) {
		t.Error("Expected decrypted message to match original message")
	}

	decryptedMessage, err = Decrypt(key2, encryptedMessage)
	if err == nil {
		t.Error("Expected decrypt message with wrong key to fail")
	}

	// DeriveKey
	passphrase := []byte("supersecret")
	salt, err := RandBytes(55)
	if err != nil {
		t.Error("Expected to create salt - ", err)
	}
	derivedKey, err := DeriveKey(passphrase, salt)
	if err != nil {
		t.Error("Expected to derive key - ", err)
	}
	if len(derivedKey) != KeySize {
		t.Error("Expected derived key length to be key size")
	}
	if bytes.Equal(derivedKey[:], passphrase) {
		t.Error("Expected derived key to be different from passphrase")
	}

	derivedKey2, err := DeriveKey(passphrase, salt)
	if err != nil {
		t.Error("Expected to derive second key - ", err)
	}
	if !bytes.Equal(derivedKey[:], derivedKey2[:]) {
		t.Error("Expected derived keys to match")
	}

	passphrase2 := []byte("notreallysecret")
	derivedKey3, err := DeriveKey(passphrase2, salt)
	if err != nil {
		t.Error("Expected to derive third key - ", err)
	}
	if bytes.Equal(derivedKey[:], derivedKey3[:]) {
		t.Error("Expected derived keys to not match")
	}

	salt2, err := RandBytes(55)
	if err != nil {
		t.Error("Expected to create second salt - ", err)
	}
	derivedKey4, err := DeriveKey(passphrase2, salt2)
	if err != nil {
		t.Error("Expected to derive fourth key - ", err)
	}
	if bytes.Equal(derivedKey3[:], derivedKey4[:]) {
		t.Error("Expected derived keys to not match")
	}

	// DeriveKeyAndSalt
	derivedKey5, salt3, err := DeriveKeyAndSalt(passphrase)
	if err != nil {
		t.Error("Expected to derive key and salt - ", err)
	}
	if len(derivedKey5[:]) != KeySize {
		t.Error("Expected key length to be key size")
	}
	if len(salt3) != SaltSize {
		t.Error("Expected salt length to be salt size")
	}
	derivedKey6, salt4, err := DeriveKeyAndSalt(passphrase)
	if err != nil {
		t.Error("Expected to derive key and salt - ", err)
	}
	if bytes.Equal(derivedKey5[:], derivedKey6[:]) {
		t.Error("Expected derived keys not to match")
	}
	if bytes.Equal(salt3, salt4) {
		t.Error("Expected salts not to match")
	}

	// DeriveSaltedKey
	saltedKey, err := DeriveSaltedKey(passphrase2)
	if err != nil {
		t.Error("Expected to derive salted key")
	}
	saltedKey2, err := DeriveSaltedKey(passphrase2)
	if err != nil {
		t.Error("Expected to derive second salted key")
	}
	if bytes.Equal(saltedKey, saltedKey2) {
		t.Error("Expected salted keys not to match")
	}

	// ExtractSalt
	extractedSalt, unsaltedKey := ExtractSalt(saltedKey)
	if len(unsaltedKey[:]) != KeySize {
		t.Error("Expected unsalted key to match key size")
	}

	rederivedKey, err := DeriveKey(passphrase2, extractedSalt)
	if err != nil {
		t.Error("Expected to derive key with extracted salt")
	}
	if !bytes.Equal(rederivedKey[:], unsaltedKey[:]) {
		t.Error("Expected rederived key to match unsalted key - ", extractedSalt, " - ", saltedKey, " - ", rederivedKey)
	}

	// Seal
	unsealedMessage := []byte("this is an unsealed message")
	sealedMessage, err := Seal(derivedKey[:], unsealedMessage)
	if err != nil {
		t.Error("Expected to seal message")
	}
	if bytes.Equal(unsealedMessage, sealedMessage) {
		t.Error("Expected sealed message to not match unsealed message")
	}

	resealedMessage, err := Seal(derivedKey[:], unsealedMessage)
	if err != nil {
		t.Error("Expected to reseal message")
	}
	if bytes.Equal(resealedMessage, sealedMessage) {
		t.Error("Expected resealed message to be different from original sealed message")
	}

	// Open
	openedMessage, err := Open(derivedKey[:], sealedMessage)
	if err != nil {
		t.Error("Expected to open message")
	}
	if !bytes.Equal(unsealedMessage, openedMessage) {
		t.Error("Expected opened message to match unsealed message")
	}

	openedMessage2, err := Open(passphrase, sealedMessage)
	if err == nil {
		t.Error("Expected open with wrong passphrase to fail")
	}
	if openedMessage2 != nil {
		t.Error("Expected failed opened message to be nil")
	}
}
