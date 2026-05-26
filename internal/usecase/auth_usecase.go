package usecase

import (
	"context"

	"REFACTORING_MAUNA/internal/dto"
)

type TokenClaims struct {
	UserID   int64
	Subject  string
	Username string
	Email    string
	Role     string
}

type TokenManager interface {
	GenerateAccessToken(userID int64, username, email, role string) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	VerifyToken(token string) (TokenClaims, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (dto.AuthResponse, error)
}
