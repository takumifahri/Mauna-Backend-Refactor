package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type JWTManager struct {
    secretKey string
}

func NewJWTManager() *JWTManager {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        secretKey = "your-secret-key-change-this" // Default dev key
    }
    return &JWTManager{secretKey: secretKey}
}

// GenerateToken generates JWT access token
func (jm *JWTManager) GenerateAccessToken(userID int64, username, email, role string) (string, error) {
    claims := JWTClaims{
        UserID:   userID,
        Username: username,
        Email:    email,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jm.secretKey))
    if err != nil {
        return "", fmt.Errorf("failed to generate token: %w", err)
    }

    return tokenString, nil
}

// GenerateRefreshToken generates JWT refresh token
func (jm *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   fmt.Sprintf("%d", userID),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jm.secretKey))
    if err != nil {
        return "", fmt.Errorf("failed to generate refresh token: %w", err)
    }

    return tokenString, nil
}

// VerifyToken verifies JWT token
func (jm *JWTManager) VerifyToken(tokenString string) (*JWTClaims, error) {
    claims := &JWTClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jm.secretKey), nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

// GetTokenExpiry returns token expiry time
func (jm *JWTManager) GetTokenExpiry(tokenString string) (time.Time, error) {
    claims, err := jm.VerifyToken(tokenString)
    if err != nil {
        return time.Time{}, err
    }

    return claims.ExpiresAt.Time, nil
}