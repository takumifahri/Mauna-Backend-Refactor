package usecase

import (
	"context"

	"REFACTORING_MAUNA/internal/dto"
)

type ProfileUsecase interface {
	GetProfile(ctx context.Context, userID int64) (dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID int64, req dto.UpdateProfileRequest) (dto.ProfileResponse, error)
	ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error
	DeactivateAccount(ctx context.Context, userID int64) error
}
