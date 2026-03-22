package security

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "fmt"
    "strings"

    "golang.org/x/crypto/argon2"
)

// Argon2Config holds argon2 parameters
type Argon2Config struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLength  uint32
    KeyLength   uint32
}

// DefaultArgon2Config returns default argon2 configuration
// Parameters based on OWASP recommendations (Dec 2023)
func DefaultArgon2Config() Argon2Config {
    return Argon2Config{
        Memory:      64 * 1024, // 64 MB - recommended minimum for interactive use
        Iterations:  3,         // 3 passes through memory
        Parallelism: 2,         // 2 parallel threads
        SaltLength:  16,        // 128-bit salt
        KeyLength:   32,        // 256-bit key (32 bytes)
    }
}

// HashPassword hash password using argon2id (most secure variant)
// Returns hash in PHC string format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
func HashPassword(password string) (string, error) {
    config := DefaultArgon2Config()

    // Generate random salt
    salt := make([]byte, config.SaltLength)
    _, err := rand.Read(salt)
    if err != nil {
        return "", fmt.Errorf("failed to generate salt: %w", err)
    }

    // Hash password with argon2id (hybrid of argon2d and argon2i)
    hash := argon2.IDKey(
        []byte(password),
        salt,
        config.Iterations,
        config.Memory,
        config.Parallelism,
        config.KeyLength,
    )

    // Encode salt and hash to base64
    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    // Format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
    hashedPassword := fmt.Sprintf(
        "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version,
        config.Memory,
        config.Iterations,
        config.Parallelism,
        b64Salt,
        b64Hash,
    )

    return hashedPassword, nil
}

// VerifyPassword verify password against argon2id hash using constant-time comparison
// Extracts config from hash string to validate with same parameters
func VerifyPassword(hashedPassword, password string) bool {
    // Parse the hash format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
    parts := strings.Split(hashedPassword, "$")
    if len(parts) != 6 {
        return false
    }

    // Verify it's argon2id
    if parts[1] != "argon2id" {
        return false
    }

    // Extract config parameters from hash
    var memory, iterations uint32
    var parallelism uint8

    _, err := fmt.Sscanf(
        parts[3],
        "m=%d,t=%d,p=%d",
        &memory,
        &iterations,
        &parallelism,
    )
    if err != nil {
        return false
    }

    // Decode salt from hash
    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return false
    }

    // Decode stored hash
    hash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false
    }

    // Hash the provided password with extracted config
    keyLength := uint32(len(hash))
    newHash := argon2.IDKey(
        []byte(password),
        salt,
        iterations,
        memory,
        parallelism,
        keyLength,
    )

    // Use constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare(hash, newHash) == 1
}

// HashPasswordWithConfig allows custom argon2 configuration
func HashPasswordWithConfig(password string, config Argon2Config) (string, error) {
    salt := make([]byte, config.SaltLength)
    _, err := rand.Read(salt)
    if err != nil {
        return "", fmt.Errorf("failed to generate salt: %w", err)
    }

    hash := argon2.IDKey(
        []byte(password),
        salt,
        config.Iterations,
        config.Memory,
        config.Parallelism,
        config.KeyLength,
    )

    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    hashedPassword := fmt.Sprintf(
        "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version,
        config.Memory,
        config.Iterations,
        config.Parallelism,
        b64Salt,
        b64Hash,
    )

    return hashedPassword, nil
}