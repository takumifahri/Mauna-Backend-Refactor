package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/internal/domain/constant"
	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/internal/usecase"
	"REFACTORING_MAUNA/pkg/security"

	"github.com/google/uuid"
)

// authService implementation
type authService struct {
	userRepo     domain.UserRepository
	tokenManager usecase.TokenManager
}

// Constructor
func NewAuthService(userRepo domain.UserRepository, tokenManager usecase.TokenManager) usecase.AuthUsecase {
	return &authService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
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
			UniqueID:   user.UniqueID,
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
		UniqueID:     uuid.NewString(),
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
		slog.ErrorContext(ctx, "register_get_created_user_failed", slog.Any("error", err), slog.Int64("user_id", userID))
		return dto.RegisterResponse{}, domain.NewInternalError(err)
	}

	var createdName string
	if createdUser.Nama != nil {
		createdName = *createdUser.Nama
	}

	return dto.RegisterResponse{
		ID:        createdUser.ID,
		UniqueID:  createdUser.UniqueID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Name:      createdName,
		Role:      string(createdUser.Role),
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

// ChangePassword implementation
func (s *authService) ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error {
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
	// Invalidate refresh token (implementation depends on how you store tokens)
	// Delete token from cookies

	return nil
}

// RefreshToken validates a refresh token and issues a new token pair.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (dto.AuthResponse, error) {
	if refreshToken == "" {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

	claims, err := s.tokenManager.VerifyToken(refreshToken)
	if err != nil {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return dto.AuthResponse{}, domain.ErrUnauthorized
	}

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
			UniqueID:   user.UniqueID,
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
