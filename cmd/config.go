package main

// DEFAULT_PWD_LENGTH is default password length
const DEFAULT_PWD_LENGTH = 8

// MAX_PWD_LENGTH is maximum password length
const MAX_PWD_LENGTH = 255 // just to limit user input

// DEFAULT_SALT_LENGTH is default length of salt used to derive the key
const DEFAULT_SALT_LENGTH = 8

// DEFAULT_NONCE_LENGTH is default length of nonce used together with the key during encryption process
const DEFAULT_NONCE_LENGTH = 12

// DEFAULT_KEY_DERIVATION_ITERATIONS is default amount of iterations for key derivation process
const DEFAULT_KEY_DERIVATION_ITERATIONS = 210000

// DEFAULT_KEY_DERIVATION_LENGTH is default length of nonce used together with the key during encryption process
const DEFAULT_KEY_DERIVATION_LENGTH = 32

type (
	// Config is the application configuration
	Config struct {
		// InputPath is a path to the input file
		InputPath string
		// OutputPath is a path to the output file
		OutputPath string
		// MinPwdLength is a minimal requirement for encryption password length
		MinPwdLength int
		// MaxPwdLength is a maximum encryption password length
		MaxPwdLength int
		// SaltLength is a length of salt used to derive the key
		SaltLength int
		// KeyDerivationIterations is a amount of iterations for key derivation process
		KeyDerivationIterations int
		// KeyDerivationLength is a length of derived key
		KeyDerivationLength int
		// NonceLength is a length of nonce used together with the key during encryption process
		NonceLength int
		// NoBase64Wrapping is a flag to determine if the wrapping/unwrapping of DTO package should be performed
		NoBase64Wrapping bool
	}
)

// NewConfig creates a new application config
func NewConfig() *Config {
	return &Config{
		MaxPwdLength: MAX_PWD_LENGTH,
	}
}
