package security

import (
    "crypto/md5"
    "crypto/sha256"
    "crypto/sha512"
    "fmt"
    "io"
)

// HashSHA256 generates SHA256 hash
func HashSHA256(data string) string {
    hash := sha256.Sum256([]byte(data))
    return fmt.Sprintf("%x", hash)
}

// HashSHA512 generates SHA512 hash
func HashSHA512(data string) string {
    hash := sha512.Sum512([]byte(data))
    return fmt.Sprintf("%x", hash)
}

// HashMD5 generates MD5 hash (use with caution)
func HashMD5(data string) string {
    hash := md5.Sum([]byte(data))
    return fmt.Sprintf("%x", hash)
}

// GenerateETag generates ETag for caching
func GenerateETag(data string) string {
    return HashSHA256(data)[:16]
}

// ComputeFileHash computes SHA256 hash of file
func ComputeFileHash(file io.Reader) (string, error) {
    hash := sha256.New()
    _, err := io.Copy(hash, file)
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}