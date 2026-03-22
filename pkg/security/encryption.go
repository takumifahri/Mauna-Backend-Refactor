package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "os"
)

type Encryptor struct {
    key []byte
}

// NewEncryptor creates new encryptor instance
func NewEncryptor() (*Encryptor, error) {
    // Get key from env or use default (32 bytes = 256-bit)
    keyStr := os.Getenv("ENCRYPTION_KEY")
    if keyStr == "" {
        keyStr = "your-32-byte-key-change-this-!!!" // Default dev key
    }

    key := []byte(keyStr)
    if len(key) != 32 {
        return nil, fmt.Errorf("encryption key must be 32 bytes, got %d", len(key))
    }

    return &Encryptor{key: key}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", fmt.Errorf("failed to generate nonce: %w", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    encoded := base64.StdEncoding.EncodeToString(ciphertext)

    return encoded, nil
}

// Decrypt decrypts encryptedData using AES-256-GCM
func (e *Encryptor) Decrypt(encryptedData string) (string, error) {
    decoded, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return "", fmt.Errorf("failed to decode: %w", err)
    }

    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", fmt.Errorf("failed to create cipher: %w", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("failed to create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()
    if len(decoded) < nonceSize {
        return "", fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("failed to decrypt: %w", err)
    }

    return string(plaintext), nil
}