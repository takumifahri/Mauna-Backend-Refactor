package usecase

import (
	"context"

	// "REFACTORING_MAUNA/internal/domain"
	// "REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/internal/dto"
	// "REFACTORING_MAUNA/pkg/security"
)

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (dto.AuthResponse, error)
}
