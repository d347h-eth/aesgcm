package codec

import (
	"github.com/d347h-eth/aesgcm/internal/domain"

	"fmt"
)

type (
	// Codec is a component responsible for encryption/decryption of data
	Codec struct {
		cfg    Config
		cipher Cipher
		rnd    RandomnessProvider
	}

	// Config ...
	Config struct {
		SaltLength              int
		KeyDerivationIterations int
		KeyDerivationLength     int
		NonceLength             int
	}

	// Cipher ...
	Cipher interface {
		Encrypt(password []byte, salt []byte, derivationIter int, derivationLength int, nonce []byte, plaintext []byte) ([]byte, error)
		Decrypt(password []byte, salt []byte, derivationIter int, derivationLength int, nonce []byte, ciphertext []byte) ([]byte, error)
	}

	// RandomnessProvider ...
	RandomnessProvider interface {
		GetRandomBytes(length int) ([]byte, error)
	}
)

func NewCodec(cfg Config, cipher Cipher, rnd RandomnessProvider) *Codec {
	return &Codec{cfg, cipher, rnd}
}

// Encrypt ...
func (c Codec) Encrypt(password []byte, plaintext []byte) (*domain.DTO, error) {
	// generate randomness
	salt, err := c.rnd.GetRandomBytes(c.cfg.SaltLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	nonce, err := c.rnd.GetRandomBytes(c.cfg.NonceLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	// encrypt
	ciphertext, err := c.cipher.Encrypt(
		password,
		salt,
		c.cfg.KeyDerivationIterations,
		c.cfg.KeyDerivationLength,
		nonce,
		plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to create ciphertext: %w", err)
	}
	// pack encrypted data into DTO type
	dto := domain.NewDTO(salt, c.cfg.KeyDerivationIterations, c.cfg.KeyDerivationLength, nonce, ciphertext)
	return dto, nil
}

// Decrypt ...
func (c Codec) Decrypt(password []byte, dto *domain.DTO) ([]byte, error) {
	plaintext, err := c.cipher.Decrypt(
		password,
		dto.KeyDerivation.Salt,
		dto.KeyDerivation.Iterations,
		dto.KeyDerivation.Length,
		dto.Nonce,
		dto.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to perform decryption: %w", err)
	}
	return plaintext, nil
}
