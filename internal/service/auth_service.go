package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"strings"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/domain/constant"
	"REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/internal/usecase"
	"REFACTORING_MAUNA/pkg/security"
)

// authService implementation
type authService struct {
	userRepo               domain.UserRepository
	tokenBlacklistRepo     domain.TokenBlacklistRepository
	passwordResetTokenRepo domain.PasswordResetTokenRepository
	passwordResetMailer    usecase.PasswordResetMailer
	tokenManager           usecase.TokenManager
	passwordResetBaseURL   string
}

// Constructor
func NewAuthService(
	userRepo domain.UserRepository,
	tokenBlacklistRepo domain.TokenBlacklistRepository,
	passwordResetTokenRepo domain.PasswordResetTokenRepository,
	passwordResetMailer usecase.PasswordResetMailer,
	tokenManager usecase.TokenManager,
) usecase.AuthUsecase {
	return &authService{
		userRepo:               userRepo,
		tokenBlacklistRepo:     tokenBlacklistRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
		passwordResetMailer:    passwordResetMailer,
		tokenManager:           tokenManager,
		passwordResetBaseURL:   os.Getenv("PASSWORD_RESET_BASE_URL"),
	}
}

// Login implementation
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	// Validate input
	if req.EmailOrUsername == "" || req.Password == "" {
		return dto.LoginResponse{}, domain.ErrInvalidCredentials
	}

	// 1. Get user from database
	user, err := s.userRepo.GetByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		return dto.LoginResponse{}, domain.ErrInvalidCredentials
	}

	// 2. Verify password
	if !security.VerifyPassword(user.PasswordHash, req.Password) {
		return dto.LoginResponse{}, domain.ErrInvalidCredentials
	}

	// 3. Check if user is active
	if !user.IsActive {
		return dto.LoginResponse{}, fmt.Errorf("account is deactivated")
	}

	// 4. Generate tokens
	accessToken, err := s.tokenManager.GenerateAccessToken(
		user.ID,
		user.Username,
		user.Email,
		string(user.Role), // Convert UserRole to string
	)
	if err != nil {
		return dto.LoginResponse{}, domain.ErrInternal
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.LoginResponse{}, domain.ErrInternal
	}

	// 5. Handle nullable Nama field
	var name string
	if user.Nama != nil {
		name = *user.Nama
	}

	// 6. Return response
	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24 hours in seconds
		User: dto.UserDataResponse{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Name:       name,              // Use dereferenced value
			Role:       string(user.Role), // Convert UserRole to string
			IsActive:   user.IsActive,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
	}, nil
}

// Register implementation
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// Validation
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Name == "" {
		return dto.RegisterResponse{}, fmt.Errorf("all fields are required")
	}

	// Check if email exists
	exists, err := s.userRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		slog.ErrorContext(ctx, "register_check_email_failed", slog.Any("error", err), slog.String("email", req.Email))
		return dto.RegisterResponse{}, domain.NewInternalError(err)
	}
	if exists {
		return dto.RegisterResponse{}, domain.ErrUserAlreadyExists
	}

	// Check if username exists
	exists, err = s.userRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		slog.ErrorContext(ctx, "register_check_username_failed", slog.Any("error", err), slog.String("username", req.Username))
		return dto.RegisterResponse{}, domain.NewInternalError(err)
	}
	if exists {
		return dto.RegisterResponse{}, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		slog.ErrorContext(ctx, "register_hash_password_failed", slog.Any("error", err))
		return dto.RegisterResponse{}, domain.NewInternalError(err)
	}

	// Convert name to pointer
	name := req.Name

	// Create user in database
	user := &entities.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Nama:         &name,             // Use pointer
		Role:         constant.RoleUser, // Use constant
		IsActive:     true,
		IsVerified:   false,
	}

	userID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "register_create_user_failed", slog.Any("error", err), slog.String("email", req.Email), slog.String("username", req.Username))
		return dto.RegisterResponse{}, domain.NewInternalError(err)
	}

	createdUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "register_get_created_user_failed", slog.Any("error", err), slog.String("user_id", userID))
		return dto.RegisterResponse{}, domain.NewInternalError(err)
	}

	var createdName string
	if createdUser.Nama != nil {
		createdName = *createdUser.Nama
	}

	return dto.RegisterResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Name:      createdName,
		Role:      string(createdUser.Role),
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	email := strings.TrimSpace(req.Email)
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return domain.ErrInvalidEmail
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Keep response generic so this endpoint cannot be used for email enumeration.
		return nil
	}
	if !user.IsActive {
		return nil
	}

	resetToken, err := security.GenerateSecureToken(32)
	if err != nil {
		slog.ErrorContext(ctx, "forgot_password_token_generate_failed", slog.Any("error", err), slog.String("user_id", user.ID))
		return domain.NewInternalError(err)
	}

	tokenHash := security.HashSHA256(resetToken)
	expiresAt := time.Now().Add(15 * time.Minute)
	if err := s.passwordResetTokenRepo.DeleteActiveByUserID(ctx, user.ID); err != nil {
		slog.ErrorContext(ctx, "forgot_password_delete_active_tokens_failed", slog.Any("error", err), slog.String("user_id", user.ID))
		return domain.NewInternalError(err)
	}

	if err := s.passwordResetTokenRepo.Create(ctx, &entities.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}); err != nil {
		slog.ErrorContext(ctx, "forgot_password_create_token_failed", slog.Any("error", err), slog.String("user_id", user.ID))
		return domain.NewInternalError(err)
	}

	if err := s.passwordResetMailer.SendPasswordReset(ctx, user.Email, authStringValue(user.Nama), resetToken, s.resetURL(resetToken)); err != nil {
		slog.ErrorContext(ctx, "forgot_password_send_email_failed", slog.Any("error", err), slog.String("user_id", user.ID), slog.String("email", user.Email))
		return domain.NewInternalError(err)
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	token := strings.TrimSpace(req.Token)
	if token == "" {
		return domain.ErrInvalidResetToken
	}
	if strings.TrimSpace(req.NewPassword) == "" || len(req.NewPassword) < 6 {
		return domain.ErrPasswordTooShort
	}

	resetToken, err := s.passwordResetTokenRepo.GetValidByHash(ctx, security.HashSHA256(token))
	if err != nil {
		if err == domain.ErrInvalidResetToken {
			return err
		}
		slog.ErrorContext(ctx, "reset_password_get_token_failed", slog.Any("error", err))
		return domain.NewInternalError(err)
	}

	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return domain.ErrInvalidResetToken
	}
	if !user.IsActive {
		return domain.ErrUnauthorized
	}

	hashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		slog.ErrorContext(ctx, "reset_password_hash_failed", slog.Any("error", err), slog.String("user_id", user.ID))
		return domain.NewInternalError(err)
	}

	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.ErrorContext(ctx, "reset_password_update_failed", slog.Any("error", err), slog.String("user_id", user.ID))
		return domain.NewInternalError(err)
	}

	if err := s.passwordResetTokenRepo.MarkUsed(ctx, resetToken.ID); err != nil {
		slog.ErrorContext(ctx, "reset_password_mark_used_failed", slog.Any("error", err), slog.String("user_id", user.ID))
		return domain.NewInternalError(err)
	}

	return nil
}

// ChangePassword implementation
func (s *authService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Verify old password
	if !security.VerifyPassword(user.PasswordHash, req.OldPassword) {
		return domain.ErrInvalidCredentials
	}

	// Hash new password
	newHashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return domain.ErrInternal
	}

	// Update password
	user.PasswordHash = newHashedPassword
	return s.userRepo.Update(ctx, user)
}

// Logout implementation
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}

	return s.revokeRefreshToken(ctx, refreshToken, "logout")
}

// RefreshToken validates a refresh token and issues a new token pair.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (dto.AuthResponse, error) {
	if refreshToken == "" {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

	blacklisted, err := s.tokenBlacklistRepo.Exists(ctx, refreshToken)
	if err != nil {
		slog.ErrorContext(ctx, "refresh_token_blacklist_check_failed", slog.Any("error", err))
		return dto.AuthResponse{}, domain.NewInternalError(err)
	}
	if blacklisted {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

	claims, err := s.tokenManager.VerifyToken(refreshToken)
	if err != nil {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

	userID := claims.Subject

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.AuthResponse{}, domain.ErrUserNotFound
	}
	if !user.IsActive {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Username, user.Email, string(user.Role))
	if err != nil {
		return dto.AuthResponse{}, domain.ErrInternal
	}

	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, domain.ErrInternal
	}

	if err := s.revokeRefreshToken(ctx, refreshToken, "rotated"); err != nil {
		return dto.AuthResponse{}, err
	}

	var name string
	if user.Nama != nil {
		name = *user.Nama
	}

	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    86400,
		User: dto.UserDataResponse{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Name:       name,
			Role:       string(user.Role),
			IsActive:   user.IsActive,
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt,
		},
	}, nil
}

func (s *authService) revokeRefreshToken(ctx context.Context, refreshToken string, reason string) error {
	claims, err := s.tokenManager.VerifyToken(refreshToken)
	if err != nil {
		return nil
	}

	userID := claims.Subject
	if userID == "" {
		return nil
	}

	expiresAt, err := s.tokenManager.GetTokenExpiry(refreshToken)
	if err != nil {
		return nil
	}

	if !expiresAt.After(time.Now()) {
		return nil
	}

	if err := s.tokenBlacklistRepo.Create(ctx, &entities.TokenBlacklist{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Reason:    &reason,
	}); err != nil {
		slog.ErrorContext(ctx, "refresh_token_revoke_failed", slog.Any("error", err), slog.String("user_id", userID), slog.String("reason", reason))
		return domain.NewInternalError(err)
	}

	return nil
}

func (s *authService) resetURL(token string) string {
	if strings.TrimSpace(s.passwordResetBaseURL) == "" {
		return ""
	}

	separator := "?"
	if strings.Contains(s.passwordResetBaseURL, "?") {
		separator = "&"
	}
	return s.passwordResetBaseURL + separator + "token=" + token
}

func authStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
