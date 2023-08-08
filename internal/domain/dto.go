package domain

import (
	"encoding/base64"
	"encoding/json"
)

type (
	// DTO contains all components used for cryptographic algorithm to encrypt/decrypt data
	DTO struct {
		KeyDerivation KeyDerivation
		Nonce         []byte
		Ciphertext    []byte
	}

	// KeyDerivation contains parameters used for key derivation algorithm
	KeyDerivation struct {
		Salt       []byte
		Iterations int
		Length     int
	}

	// DTOBase64 is a proxy type that represents DTO with Base64 encoding
	DTOBase64 struct {
		KeyDerivation KeyDerivationBase64 `json:"key_derivation"`
		Nonce         string              `json:"nonce"`
		Ciphertext    string              `json:"ciphertext"`
	}

	// KeyDerivationBase64 ...
	KeyDerivationBase64 struct {
		Salt       string `json:"salt"`
		Iterations int    `json:"iterations"`
		Length     int    `json:"length"`
	}
)

// NewDTO creates a new instance of DTO
func NewDTO(salt []byte, keyDerivationIterations int, keyDerivationLength int, nonce []byte, ciphertext []byte) *DTO {
	dto := DTO{
		KeyDerivation: KeyDerivation{
			Salt:       salt,
			Iterations: keyDerivationIterations,
			Length:     keyDerivationLength,
		},
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}
	return &dto
}

// MarshalJSON implements JSON marshaling into Base64 representation
func (m DTO) MarshalJSON() ([]byte, error) {
	tempStruct := DTOBase64{
		KeyDerivation: KeyDerivationBase64{
			Salt:       base64.StdEncoding.EncodeToString(m.KeyDerivation.Salt),
			Iterations: m.KeyDerivation.Iterations,
			Length:     m.KeyDerivation.Length,
		},
		Nonce:      base64.StdEncoding.EncodeToString(m.Nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(m.Ciphertext),
	}
	return json.Marshal(tempStruct)
}

// UnmarshalJSON implements JSON unmarshaling from Base64 representation
func (m *DTO) UnmarshalJSON(data []byte) error {
	tempStruct := DTOBase64{}
	err := json.Unmarshal(data, &tempStruct)
	if err != nil {
		return err
	}
	m.KeyDerivation.Salt, err = base64.StdEncoding.DecodeString(tempStruct.KeyDerivation.Salt)
	if err != nil {
		return err
	}
	m.KeyDerivation.Iterations = tempStruct.KeyDerivation.Iterations
	m.KeyDerivation.Length = tempStruct.KeyDerivation.Length
	m.Nonce, err = base64.StdEncoding.DecodeString(tempStruct.Nonce)
	if err != nil {
		return err
	}
	m.Ciphertext, err = base64.StdEncoding.DecodeString(tempStruct.Ciphertext)
	if err != nil {
		return err
	}
	return nil
}
