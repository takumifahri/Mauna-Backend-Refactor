package service

import (
	"context"
	"net/mail"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/internal/usecase"
	"REFACTORING_MAUNA/pkg/security"
)

type managementUsersService struct {
	userRepo domain.ManagementUsersRepository
}

var _ usecase.ManagementUsersUsecase = (*managementUsersService)(nil)

func NewManagementUsersService(userRepo domain.ManagementUsersRepository) usecase.ManagementUsersUsecase {
	return &managementUsersService{userRepo: userRepo}
}

func (s *managementUsersService) ListUsers(ctx context.Context, filter admin.ManagementUsersFilter) (admin.ManagementUsersListResponse, error) {
	return s.userRepo.List(ctx, filter)
}

func (s *managementUsersService) GetUser(ctx context.Context, id string) (admin.ManagementUserResponse, error) {
	return s.userRepo.GetByID(ctx, id, false)
}

func (s *managementUsersService) CreateUser(ctx context.Context, req admin.CreateManagementUserRequest) (admin.ManagementUserResponse, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.TrimSpace(req.Email)
	name := strings.TrimSpace(req.Name)
	phone := strings.TrimSpace(req.Phone)
	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = "user"
	}
	role = strings.ToLower(role)

	if err := validateManagementUserInput(username, email, req.Password, name, role); err != nil {
		return admin.ManagementUserResponse{}, err
	}

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return admin.ManagementUserResponse{}, domain.NewInternalError(err)
	}

	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	isVerified := false
	if req.IsVerified != nil {
		isVerified = *req.IsVerified
	}

	createReq := admin.CreateManagementUserRequest{
		Username:   username,
		Email:      email,
		Password:   passwordHash,
		Name:       name,
		Phone:      phone,
		Role:       role,
		Avatar:     strings.TrimSpace(req.Avatar),
		AvatarURL:  strings.TrimSpace(req.AvatarURL),
		IsActive:   &isActive,
		IsVerified: &isVerified,
	}

	id, err := s.userRepo.Create(ctx, createReq)
	if err != nil {
		return admin.ManagementUserResponse{}, err
	}
	return s.userRepo.GetByID(ctx, id, false)
}

func (s *managementUsersService) UpdateUser(ctx context.Context, id string, req admin.UpdateManagementUserRequest) (admin.ManagementUserResponse, error) {
	if _, err := s.userRepo.GetByID(ctx, id, false); err != nil {
		return admin.ManagementUserResponse{}, err
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if len(username) < 3 {
			return admin.ManagementUserResponse{}, domain.NewBusinessError("INVALID_USERNAME", "username must be at least 3 characters", domain.ErrInvalidRequest)
		}
		req.Username = &username
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if err := validateEmail(email); err != nil {
			return admin.ManagementUserResponse{}, err
		}
		req.Email = &email
	}
	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if password == "" {
			req.Password = nil
		} else {
			if len(password) < 6 {
				return admin.ManagementUserResponse{}, domain.ErrPasswordTooShort
			}
			passwordHash, err := security.HashPassword(password)
			if err != nil {
				return admin.ManagementUserResponse{}, domain.NewInternalError(err)
			}
			req.Password = &passwordHash
		}
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return admin.ManagementUserResponse{}, domain.NewBusinessError("INVALID_NAME", "name cannot be empty", domain.ErrInvalidRequest)
		}
		req.Name = &name
	}
	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		req.Phone = &phone
	}
	if req.Avatar != nil {
		avatar := strings.TrimSpace(*req.Avatar)
		req.Avatar = &avatar
	}
	if req.AvatarURL != nil {
		avatarURL := strings.TrimSpace(*req.AvatarURL)
		req.AvatarURL = &avatarURL
	}

	if req.Role != nil {
		role := strings.ToLower(strings.TrimSpace(*req.Role))
		if !isValidManagementUserRole(role) {
			return admin.ManagementUserResponse{}, domain.NewBusinessError("INVALID_ROLE", "role must be admin, user, or moderator", domain.ErrInvalidRequest)
		}
		req.Role = &role
	}

	if err := s.userRepo.Update(ctx, id, req); err != nil {
		return admin.ManagementUserResponse{}, err
	}
	return s.userRepo.GetByID(ctx, id, false)
}

func (s *managementUsersService) SoftDeleteUser(ctx context.Context, id string) error {
	return s.userRepo.SoftDelete(ctx, id)
}

func (s *managementUsersService) HardDeleteUser(ctx context.Context, id string) error {
	return s.userRepo.HardDelete(ctx, id)
}

func (s *managementUsersService) RestoreUser(ctx context.Context, id string) (admin.ManagementUserResponse, error) {
	if err := s.userRepo.Restore(ctx, id); err != nil {
		return admin.ManagementUserResponse{}, err
	}
	return s.userRepo.GetByID(ctx, id, false)
}

func validateManagementUserInput(username, email, password, name, role string) error {
	if len(username) < 3 {
		return domain.NewBusinessError("INVALID_USERNAME", "username must be at least 3 characters", domain.ErrInvalidRequest)
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if len(strings.TrimSpace(password)) < 6 {
		return domain.ErrPasswordTooShort
	}
	if strings.TrimSpace(name) == "" {
		return domain.NewBusinessError("INVALID_NAME", "name cannot be empty", domain.ErrInvalidRequest)
	}
	if !isValidManagementUserRole(role) {
		return domain.NewBusinessError("INVALID_ROLE", "role must be admin, user, or moderator", domain.ErrInvalidRequest)
	}
	return nil
}

func validateEmail(email string) error {
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return domain.ErrInvalidEmail
	}
	return nil
}

func isValidManagementUserRole(role string) bool {
	switch role {
	case "admin", "user", "moderator":
		return true
	default:
		return false
	}
}
