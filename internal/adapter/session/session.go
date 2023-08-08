package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/d347h-eth/aesgcm/internal/domain"
)

type (
	// Session is a component responsible for driving the core use case and user interaction
	Session struct {
		cfg      Config
		terminal Terminal
		storage  Storage
		codec    Codec
	}

	// Config ...
	Config struct {
		NoBase64Wrapping bool
	}

	// Terminal is a component responsible for receiving secrets from the user in real-time
	Terminal interface {
		ReceiveEncryptionPwd() ([]byte, error)
		ReceiveDecryptionPwd() ([]byte, error)
	}

	// Codec is a component responsible for encryption/decryption of data
	Codec interface {
		Encrypt(password []byte, plaintext []byte) (dto *domain.DTO, err error)
		Decrypt(password []byte, dto *domain.DTO) (plaintext []byte, err error)
	}

	// Storage is responsible for reading and writing of data
	Storage interface {
		ResourceExist(path string) bool
		Read(path string) ([]byte, error)
		Write(path string, data []byte) error
	}
)

// NewSession ...
func NewSession(cfg Config, terminal Terminal, storage Storage, codec Codec) *Session {
	return &Session{cfg, terminal, storage, codec}
}

// Encrypt ...
func (s Session) Encrypt(inputPath string, outputPath string) error {
	// make sure the file with input plaintext exists
	if !s.storage.ResourceExist(inputPath) {
		return fmt.Errorf("the file with plaintext input has not been found at: %q", inputPath)
	}
	// make sure the file with output ciphertext doesn't exist
	if s.storage.ResourceExist(outputPath) {
		return fmt.Errorf("the file with ciphertext already exists at %q: specify different output path with -o flag or remove the file", outputPath)
	}
	// receive input plaintext
	plaintext, err := s.storage.Read(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}
	// receive the password used to derive the key
	password, err := s.terminal.ReceiveEncryptionPwd()
	if err != nil {
		return fmt.Errorf("failed to receive a password: %w", err)
	}
	// encrypt
	dto, err := s.codec.Encrypt(password, plaintext)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}
	// output data
	outputJSON, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("an error occurred while encoding to JSON: %w", err)
	}
	var outputData []byte
	if s.cfg.NoBase64Wrapping { // output JSON without Base64 encoding
		outputData = outputJSON
	} else { // wrap output JSON with additional Base64 encoding
		tmpBuffer := make([]byte, base64.StdEncoding.EncodedLen(len(outputJSON)))
		base64.StdEncoding.Encode(tmpBuffer, outputJSON)
		outputData = tmpBuffer
	}
	err = s.storage.Write(outputPath, outputData)
	if err != nil {
		return fmt.Errorf("failed to save encrypted data: %w", err)
	}
	fmt.Printf("Successfully encrypted to %q\n", outputPath)
	return nil
}

// Decrypt ...
func (s Session) Decrypt(inputPath string, outputPath string) error {
	// make sure the file with input ciphertext exists
	if !s.storage.ResourceExist(inputPath) {
		return fmt.Errorf("the file with ciphertext input has not been found at: %q", inputPath)
	}
	// make sure the file with output plaintext doesn't exist
	if s.storage.ResourceExist(outputPath) {
		return fmt.Errorf("the file with plaintext already exists at %q: specify different output path with -o flag or remove the file",
			outputPath)
	}
	// process the input
	inputData, err := s.storage.Read(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}
	var inputJson []byte
	if s.cfg.NoBase64Wrapping { // input is already JSON
		inputJson = inputData
	} else { // decode Base64 wrapping first
		tmpBuffer := make([]byte, base64.StdEncoding.DecodedLen(len(inputData)))
		bytesWritten, err := base64.StdEncoding.Decode(tmpBuffer, inputData)
		if err != nil {
			return fmt.Errorf("failed to decode Base64 data from input file: %w", err)
		}
		inputJson = tmpBuffer[:bytesWritten]
	}
	// unmarshal JSON data into DTO type
	dto := &domain.DTO{}
	err = json.Unmarshal(inputJson, dto)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON from input file: %w", err)
	}
	// receive the password used to derive the key
	password, err := s.terminal.ReceiveDecryptionPwd()
	if err != nil {
		return fmt.Errorf("failed to receive a password: %w", err)
	}
	// decrypt
	plaintext, err := s.codec.Decrypt(password, dto)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}
	// output plaintext
	err = s.storage.Write(outputPath, plaintext)
	if err != nil {
		return fmt.Errorf("failed to save plaintext output: %w", err)
	}
	fmt.Printf("Successfully decrypted to %q\n", outputPath)
	return nil
}
