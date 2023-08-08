package terminal

import (
	"bytes"
	"fmt"

	"golang.org/x/term"
)

type (
	Terminal struct {
		cfg Config
	}

	Config struct {
		MinPwdLength int
		MaxPwdLength int
	}
)

func NewTerminal(cfg Config) *Terminal {
	return &Terminal{cfg}
}

// ReceiveEncryptionPwd promts user to enter a password for encryption (twice)
func (t Terminal) ReceiveEncryptionPwd() ([]byte, error) {
	fmt.Printf("Please enter your password (min. %d characters): ", t.cfg.MinPwdLength)
	secret, err := term.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while reading the password: %w", err)
	}
	if len(secret) < t.cfg.MinPwdLength {
		return nil, fmt.Errorf("invalid input: password must be at least %d characters", t.cfg.MinPwdLength)
	}
	if len(secret) > t.cfg.MaxPwdLength {
		return nil, fmt.Errorf("passwords longer than %d are not supported", t.cfg.MaxPwdLength)
	}
	fmt.Print("\nPlease re-enter your password: ")
	secretRepeat, err := term.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while reading the repeated password: %w", err)
	}
	if !bytes.Equal(secret, secretRepeat) {
		return nil, fmt.Errorf("passwords do not match")
	}
	fmt.Println()
	return secret, nil
}

// ReceiveDecryptionPwd promts user to enter a password for decryption
func (t Terminal) ReceiveDecryptionPwd() ([]byte, error) {
	fmt.Printf("Please enter your password: ")
	secret, err := term.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while reading the password: %w", err)
	}
	if len(secret) <= 0 {
		return nil, fmt.Errorf("decryption password can't be empty")
	}
	if len(secret) > t.cfg.MaxPwdLength {
		return nil, fmt.Errorf("passwords longer than %d are not supported", t.cfg.MaxPwdLength)
	}
	fmt.Println()
	return secret, nil
}
