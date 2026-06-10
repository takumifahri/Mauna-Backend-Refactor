package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		byteLength = 32
	}

	token := make([]byte, byteLength)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}

func GenerateNumericOTP(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	otp := make([]byte, length)
	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate otp: %w", err)
		}
		otp[i] = byte('0' + n.Int64())
	}
	return string(otp), nil
}
