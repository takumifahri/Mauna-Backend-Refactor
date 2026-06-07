package usecase

import (
	"context"
	"time"

	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/pkg/security"
)

type TokenClaims = security.TokenClaims

type TokenManager interface {
	GenerateAccessToken(userID string, username, email, role string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	VerifyToken(token string) (TokenClaims, error)
	GetTokenExpiry(token string) (time.Time, error)
}

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (dto.AuthResponse, error)
}

type PasswordResetMailer interface {
	SendPasswordReset(ctx context.Context, to string, name string, resetToken string, resetURL string) error
}
