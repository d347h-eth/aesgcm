package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"

	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type (
	AESGCM struct{}
)

func NewAESGCM() *AESGCM {
	return &AESGCM{}
}

// Encrypt performs the encryption and returns ciphertext
func (aes AESGCM) Encrypt(
	password []byte,
	salt []byte,
	derivationIter int,
	derivationLength int,
	nonce []byte,
	plaintext []byte,
) ([]byte, error) {
	key := deriveKey(password, salt, derivationIter, derivationLength)
	aesgcm, err := initCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt: %w", err)
	}
	if len(nonce) != aesgcm.NonceSize() {
		return nil, fmt.Errorf(
			"incorrect nonce size: %d, must be %d",
			len(nonce),
			aesgcm.NonceSize(),
		)
	}
	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}

// Decrypt performs the decryption and returns the plaintext
func (aes AESGCM) Decrypt(password []byte,
	salt []byte,
	derivationIter int,
	derivationLength int,
	nonce []byte,
	ciphertext []byte,
) ([]byte, error) {
	key := deriveKey(password, salt, derivationIter, derivationLength)
	aesgcm, err := initCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt %w", err)
	}
	if len(nonce) != aesgcm.NonceSize() {
		return nil, fmt.Errorf(
			"incorrect nonce size: %d, must be %d",
			len(nonce),
			aesgcm.NonceSize(),
		)
	}
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

func deriveKey(password []byte, salt []byte, derivationIter int, derivationLength int) []byte {
	return pbkdf2.Key(password, salt, derivationIter, derivationLength, sha512.New)
}

func initCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to init new cipher: %w", err)
	}
	return cipher.NewGCM(block)
}
