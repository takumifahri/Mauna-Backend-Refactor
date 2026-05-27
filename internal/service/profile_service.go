package service

import (
	"context"
	"log/slog"
	"net/mail"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/internal/usecase"
	"REFACTORING_MAUNA/pkg/security"
)

type profileService struct {
	userRepo domain.UserRepository
}

var _ usecase.ProfileUsecase = (*profileService)(nil)

func NewProfileService(userRepo domain.UserRepository) usecase.ProfileUsecase {
	return &profileService{
		userRepo: userRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context, userID int64) (dto.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "profile_get_user_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return dto.ProfileResponse{}, domain.ErrUserNotFound
	}

	if !user.IsActive {
		return dto.ProfileResponse{}, domain.ErrUnauthorized
	}

	return profileResponseFromUser(user), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID int64, req dto.UpdateProfileRequest) (dto.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "profile_update_user_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return dto.ProfileResponse{}, domain.ErrUserNotFound
	}

	if !user.IsActive {
		return dto.ProfileResponse{}, domain.ErrUnauthorized
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return dto.ProfileResponse{}, domain.NewBusinessError("INVALID_USERNAME", "username cannot be empty", domain.ErrInvalidRequest)
		}
		if username != user.Username {
			exists, err := s.userRepo.CheckUsernameExists(ctx, username)
			if err != nil {
				slog.ErrorContext(ctx, "profile_update_check_username_failed", slog.Any("error", err), slog.String("username", username))
				return dto.ProfileResponse{}, domain.NewInternalError(err)
			}
			if exists {
				return dto.ProfileResponse{}, domain.ErrUserAlreadyExists
			}
		}
		user.Username = username
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		address, err := mail.ParseAddress(email)
		if err != nil || address.Address != email {
			return dto.ProfileResponse{}, domain.ErrInvalidEmail
		}
		if email != user.Email {
			exists, err := s.userRepo.CheckEmailExists(ctx, email)
			if err != nil {
				slog.ErrorContext(ctx, "profile_update_check_email_failed", slog.Any("error", err), slog.String("email", email))
				return dto.ProfileResponse{}, domain.NewInternalError(err)
			}
			if exists {
				return dto.ProfileResponse{}, domain.ErrUserAlreadyExists
			}
		}
		user.Email = email
	}

	if req.Name != nil {
		user.Nama = trimmedStringPtr(*req.Name)
	}
	if req.Phone != nil {
		user.Telpon = trimmedStringPtr(*req.Phone)
	}
	if req.AvatarURL != nil {
		user.AvatarURL = trimmedStringPtr(*req.AvatarURL)
	}
	if req.Bio != nil {
		user.Bio = trimmedStringPtr(*req.Bio)
	}

	if err := s.userRepo.UpdateProfile(ctx, user); err != nil {
		slog.ErrorContext(ctx, "profile_update_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return dto.ProfileResponse{}, domain.NewInternalError(err)
	}

	updatedUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "profile_update_get_user_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return dto.ProfileResponse{}, domain.NewInternalError(err)
	}

	return profileResponseFromUser(updatedUser), nil
}

func (s *profileService) DeactivateAccount(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "profile_deactivate_user_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return domain.ErrUserNotFound
	}
	if !user.IsActive {
		return domain.ErrUnauthorized
	}

	if err := s.userRepo.Deactivate(ctx, userID); err != nil {
		slog.ErrorContext(ctx, "profile_deactivate_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return domain.NewInternalError(err)
	}

	return nil
}

func (s *profileService) ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, "profile_change_password_user_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return domain.ErrUserNotFound
	}

	if !user.IsActive {
		return domain.ErrUnauthorized
	}

	if strings.TrimSpace(req.NewPassword) == "" || len(req.NewPassword) < 6 {
		return domain.ErrPasswordTooShort
	}

	if !security.VerifyPassword(user.PasswordHash, req.OldPassword) {
		return domain.ErrInvalidCredentials
	}

	hashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		slog.ErrorContext(ctx, "profile_change_password_hash_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return domain.NewInternalError(err)
	}

	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.ErrorContext(ctx, "profile_change_password_update_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return domain.NewInternalError(err)
	}

	return nil
}

func profileResponseFromUser(user *entities.User) dto.ProfileResponse {
	return dto.ProfileResponse{
		ID:                    user.ID,
		UniqueID:              user.UniqueID,
		Username:              user.Username,
		Email:                 user.Email,
		Name:                  stringValue(user.Nama),
		Phone:                 stringValue(user.Telpon),
		Role:                  string(user.Role),
		IsActive:              user.IsActive,
		IsVerified:            user.IsVerified,
		Avatar:                stringValue(user.AvatarURL),
		Bio:                   stringValue(user.Bio),
		TotalBadges:           user.TotalBadges,
		CurrentStreak:         user.CurrentStreak,
		LongestStreak:         user.LongestStreak,
		WeeklyXP:              user.WeeklyXp,
		Tier:                  string(user.Tier),
		TotalXP:               user.TotalXp,
		TotalQuizzesCompleted: user.TotalQuizzesCompleted,
		TotalPoints:           user.TotalPoints,
		CreatedAt:             user.CreatedAt,
		LastLogin:             user.LastLogin,
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func trimmedStringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
