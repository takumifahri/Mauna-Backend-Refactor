package usecase

import (
	"context"

	"REFACTORING_MAUNA/internal/dto"
)

type ProfileUsecase interface {
	GetProfile(ctx context.Context, userID string) (dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (dto.ProfileResponse, error)
	ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error
	DeactivateAccount(ctx context.Context, userID string) error
}
