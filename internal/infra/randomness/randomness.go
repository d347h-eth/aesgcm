package randomness

import (
	"crypto/rand"
	"io"
)

type (
	OSRandomness struct{}
)

func NewOSRandomness() *OSRandomness {
	return &OSRandomness{}
}

// GetRandomBytes returns a slice of random bytes of specified length
func (os OSRandomness) GetRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return bytes, err
	}
	return bytes, nil
}
