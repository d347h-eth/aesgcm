package qrencoder

import (
	"github.com/skip2/go-qrcode"
)

type (
	// QREncoderConfig ...
	QREncoderConfig struct {
		// Level is error recovery level
		Level int
		// Size is generated image size in pixels
		Size int
	}

	// QREncoder ...
	QREncoder struct {
		cfg QREncoderConfig
	}
)

// NewQREncoder ...
func NewQREncoder(cfg QREncoderConfig) *QREncoder {
	return &QREncoder{cfg}
}

// Encode generates QR code image
func (e *QREncoder) Encode(data []byte) ([]byte, error) {
	var imgBytes []byte
	imgBytes, err := qrcode.Encode(
		string(data),
		qrcode.RecoveryLevel(e.cfg.Level),
		e.cfg.Size,
	)
	return imgBytes, err
}
