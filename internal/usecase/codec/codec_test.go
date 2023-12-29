package codec

import (
	"bytes"
	"testing"

	"github.com/d347h-eth/aesgcm/internal/infra/aesgcm"
	"github.com/d347h-eth/aesgcm/internal/infra/randomness"
)

func TestCodec(t *testing.T) {
	aesgcm := aesgcm.NewAESGCM()
	osRandomness := randomness.NewOSRandomness()
	cfg := Config{
		SaltLength:              128,
		NonceLength:             12,
		KeyDerivationIterations: 1000000,
		KeyDerivationLength:     32,
	}
	codec := NewCodec(cfg, aesgcm, osRandomness)
	var testInput = struct {
		pwd       []byte
		plaintext []byte
	}{
		[]byte("testpassword"),
		[]byte("secretplaintext"),
	}
	dto, err := codec.Encrypt(testInput.pwd, testInput.plaintext)
	if err != nil {
		t.Fatalf("failed encryption: %s", err)
	}
	plaintext, err := codec.Decrypt(testInput.pwd, dto)
	if err != nil {
		t.Fatalf("failed decryption: %s", err)
	}
	if !bytes.Equal(testInput.plaintext, plaintext) {
		t.Fatalf(
			"decrypted plaintext and original secret don't match: got %q, want %q",
			plaintext,
			testInput.plaintext,
		)
	}
}
