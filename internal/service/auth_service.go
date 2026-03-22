package service

import (
    "context"
    "fmt"
    "time"

    "REFACTORING_MAUNA/internal/domain"
    "REFACTORING_MAUNA/internal/domain/entities"
    "REFACTORING_MAUNA/internal/dto"
    "REFACTORING_MAUNA/pkg/security"
)

// AuthUsecase interface
type AuthUsecase interface {
    Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
    Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
    ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error
}

// authService implementation
type authService struct {
    userRepo   domain.UserRepository
    jwtManager *security.JWTManager
}

// Constructor
func NewAuthService(userRepo domain.UserRepository) AuthUsecase {
    return &authService{
        userRepo:   userRepo,
        jwtManager: security.NewJWTManager(),
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
    accessToken, err := s.jwtManager.GenerateAccessToken(
        user.ID,
        user.Username,
        user.Email,
        string(user.Role), // Convert UserRole to string
    )
    if err != nil {
        return dto.LoginResponse{}, domain.ErrInternal
    }

    refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
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
            ID:        user.ID,
            UniqueID:  user.UniqueID,
            Username:  user.Username,
            Email:     user.Email,
            Name:      name, // Use dereferenced value
            Role:      string(user.Role), // Convert UserRole to string
            IsActive:  user.IsActive,
            IsVerified: user.IsVerified,
            CreatedAt: user.CreatedAt,
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
        return dto.RegisterResponse{}, domain.ErrInternal
    }
    if exists {
        return dto.RegisterResponse{}, domain.ErrUserAlreadyExists
    }

    // Check if username exists
    exists, err = s.userRepo.CheckUsernameExists(ctx, req.Username)
    if err != nil {
        return dto.RegisterResponse{}, domain.ErrInternal
    }
    if exists {
        return dto.RegisterResponse{}, domain.ErrUserAlreadyExists
    }

    // Hash password
    hashedPassword, err := security.HashPassword(req.Password)
    if err != nil {
        return dto.RegisterResponse{}, domain.ErrInternal
    }

    // Convert name to pointer
    name := req.Name

    // Create user in database
    user := &entities.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        Nama:         &name, // Use pointer
        Role:         entities.RoleUser, // Use constant
        IsActive:     true,
        IsVerified:   false,
    }

    userID, err := s.userRepo.Create(ctx, user)
    if err != nil {
        return dto.RegisterResponse{}, domain.ErrInternal
    }

    return dto.RegisterResponse{
        ID:        userID,
        Username:  req.Username,
        Email:     req.Email,
        CreatedAt: time.Now(),
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
	return nil
}