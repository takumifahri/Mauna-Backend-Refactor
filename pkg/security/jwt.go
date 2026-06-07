package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenClaims struct {
	UserID   string
	Subject  string
	Username string
	Email    string
	Role     string
}

type JWTManager struct {
	secretKey string
}

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
	AccessTokenExpiresIn = int(AccessTokenDuration / time.Second)
)

func NewJWTManager() *JWTManager {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-this" // Default dev key
	}
	return &JWTManager{secretKey: secretKey}
}

// GenerateToken generates JWT access token
func (jm *JWTManager) GenerateAccessToken(userID string, username, email, role string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
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
func (jm *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenDuration)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jm.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken verifies JWT token
func (jm *JWTManager) VerifyToken(tokenString string) (TokenClaims, error) {
	claims, err := jm.parseJWTClaims(tokenString)
	if err != nil {
		return TokenClaims{}, err
	}

	return TokenClaims{
		UserID:   claims.UserID,
		Subject:  claims.Subject,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role,
	}, nil
}

// GetTokenExpiry returns token expiry time
func (jm *JWTManager) GetTokenExpiry(tokenString string) (time.Time, error) {
	claims, err := jm.parseJWTClaims(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	if claims.ExpiresAt == nil {
		return time.Time{}, fmt.Errorf("token has no expiry")
	}

	return claims.ExpiresAt.Time, nil
}

func (jm *JWTManager) parseJWTClaims(tokenString string) (*JWTClaims, error) {
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
