package main

import (
	"fmt"
	"os"

	"github.com/d347h-eth/aesgcm/internal/adapter/session"
	"github.com/d347h-eth/aesgcm/internal/infra/aesgcm"
	"github.com/d347h-eth/aesgcm/internal/infra/filesystem"
	"github.com/d347h-eth/aesgcm/internal/infra/qrencoder"
	"github.com/d347h-eth/aesgcm/internal/infra/randomness"
	"github.com/d347h-eth/aesgcm/internal/infra/terminal"
	"github.com/d347h-eth/aesgcm/internal/usecase/codec"

	"github.com/spf13/cobra"
)

var QRRecoveryLevels = map[string]int{
	"low":     0,
	"medium":  1,
	"high":    2,
	"highest": 3,
}

func main() {
	cfg := NewConfig()

	var cmd = &cobra.Command{
		Use:   "aesgcm [encrypt|decrypt] INPUT_FILEPATH",
		Short: "Encrypts or decrypts a file with password using AES-GCM algorithm",
		Long: `Encrypts or decrypts a file with password using AES-GCM algorithm.

Usage examples:
  aesgcm encrypt example.txt
  aesgcm decrypt example.aes`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			action := args[0]
			cfg.InputPath = args[1]
			cmd.SilenceUsage = true
			return runApp(
				action,
				cfg.InputPath,
				cfg.OutputPath,
				mapSessionCfg(cfg),
				mapCodecCfg(cfg),
				mapQREncoderCfg(cfg),
				mapTerminalCfg(cfg),
			)
		},
	}

	cmd.Flags().StringVarP(&cfg.OutputPath, "output", "o", "",
		fmt.Sprintf("Redirect output into the specified file. "+
			"By default the output is saved at %q for encryption and %q for decryption.",
			"INPUT_FILEPATH.aes", "INPUT_FILEPATH.txt"))

	cmd.Flags().IntVarP(&cfg.MinPwdLength, "min-password-length", "p", DEFAULT_PWD_LENGTH,
		"Minimum password length requirement. The password must be at least this many characters long.")

	cmd.Flags().IntVar(&cfg.SaltLength, "salt-length", DEFAULT_SALT_LENGTH,
		"Salt length. Salt is used to derive key from the password. Don't change unless you're absolutely confident.")
	cmd.Flags().IntVar(&cfg.NonceLength, "nonce-length", DEFAULT_NONCE_LENGTH,
		"Nonce length. Nonce is used together with the key to encrypt data and must be unique per given key. "+
			"Don't change unless you're absolutely confident.")
	cmd.Flags().IntVar(&cfg.KeyDerivationIterations, "key-derivation-iterations", DEFAULT_KEY_DERIVATION_ITERATIONS,
		"Amount of iterations used to derive the key. Don't change unless you're absolutely confident.")
	cmd.Flags().IntVar(&cfg.KeyDerivationLength, "key-derivation-length", DEFAULT_KEY_DERIVATION_LENGTH,
		"Length of derived key. Don't change unless you're absolutely confident.")

	cmd.Flags().BoolVar(&cfg.DisableBase64Processing, "disable-base64", false,
		"Use Base64 encoding/decoding with the input/output DTO.")

	cmd.Flags().BoolVar(&cfg.EnableQRGeneration, "qr-enable", false,
		"Generate a PNG image with QR code alongside the encoded output.")
	cmd.Flags().StringVar(&cfg.QRRecoveryLevel, "qr-level", DEFAULT_QR_RECOVERY_LEVEL,
		"Set QR code error recovery level (low, medium, high, highest).")
	cmd.Flags().IntVar(&cfg.QRSize, "qr-size", DEFAULT_QR_SIZE,
		"Set QR code image size in pixels (if positive value provided). "+
			"A negative value causes a variable sized image to be rendered "+
			"(specified value will be applied as width in pixels for each QR code \"module\").")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if cfg.EnableQRGeneration {
			return validateQRRecoveryLevel(cfg.QRRecoveryLevel)
		}
		return nil
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runApp(
	action string,
	inputPath string,
	outputPath string,
	sessionCfg session.Config,
	codecCfg codec.Config,
	qrEncoderCfg qrencoder.QREncoderConfig,
	terminalCfg terminal.Config) error {
	terminal := terminal.NewTerminal(terminalCfg)
	storage := filesystem.NewFileSystem()
	aesgcm := aesgcm.NewAESGCM()
	osRandomness := randomness.NewOSRandomness()
	codec := codec.NewCodec(codecCfg, aesgcm, osRandomness)
	qrEncoder := qrencoder.NewQREncoder(qrEncoderCfg)
	session := session.NewSession(
		sessionCfg,
		terminal,
		storage,
		codec,
		qrEncoder,
	)

	switch action {
	case "encrypt":
		if outputPath == "" {
			outputPath = inputPath + ".aes"
		}
		return session.Encrypt(inputPath, outputPath)
	case "decrypt":
		if outputPath == "" {
			outputPath = inputPath + ".txt"
		}
		return session.Decrypt(inputPath, outputPath)
	default:
		return fmt.Errorf("invalid action %q: "+
			"please choose %q or %q", action, "encrypt", "decrypt")
	}
}

func mapSessionCfg(cfg *Config) session.Config {
	return session.Config{
		Base64WrappingDisabled: cfg.DisableBase64Processing,
		QRGenerationEnabled:    cfg.EnableQRGeneration,
	}
}

func mapTerminalCfg(cfg *Config) terminal.Config {
	return terminal.Config{
		MinPwdLength: cfg.MinPwdLength,
		MaxPwdLength: cfg.MaxPwdLength,
	}
}

func mapCodecCfg(cfg *Config) codec.Config {
	return codec.Config{
		SaltLength:              cfg.SaltLength,
		KeyDerivationIterations: cfg.KeyDerivationIterations,
		KeyDerivationLength:     cfg.KeyDerivationLength,
		NonceLength:             cfg.NonceLength,
	}
}

func validateQRRecoveryLevel(level string) error {
	if _, ok := QRRecoveryLevels[level]; !ok {
		return fmt.Errorf("invalid QR code error recovery level provided: %s", level)
	}
	return nil
}

func mapQREncoderCfg(cfg *Config) qrencoder.QREncoderConfig {
	return qrencoder.QREncoderConfig{
		Level: QRRecoveryLevels[cfg.QRRecoveryLevel],
		Size:  cfg.QRSize,
	}
}
