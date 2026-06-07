package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
