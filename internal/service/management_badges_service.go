package service

import (
	"context"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type managementBadgesService struct {
	badgeRepo domain.ManagementBadgesRepository
}

var _ usecase.ManagementBadgesUsecase = (*managementBadgesService)(nil)

func NewManagementBadgesService(badgeRepo domain.ManagementBadgesRepository) usecase.ManagementBadgesUsecase {
	return &managementBadgesService{badgeRepo: badgeRepo}
}

func (s *managementBadgesService) ListBadges(ctx context.Context, filter admin.ManagementBadgesFilter) (admin.ManagementBadgesListResponse, error) {
	filter.Level = strings.ToLower(strings.TrimSpace(filter.Level))
	if filter.Level != "" && !isValidBadgeLevel(filter.Level) {
		return admin.ManagementBadgesListResponse{}, invalidBadgeLevelError()
	}
	return s.badgeRepo.List(ctx, filter)
}

func (s *managementBadgesService) GetBadge(ctx context.Context, id int64) (admin.ManagementBadgeResponse, error) {
	return s.badgeRepo.GetByID(ctx, id, false)
}

func (s *managementBadgesService) CreateBadge(ctx context.Context, req admin.CreateManagementBadgeRequest) (admin.ManagementBadgeResponse, error) {
	name := strings.TrimSpace(req.Name)
	description := strings.TrimSpace(req.Description)
	icon := strings.TrimSpace(req.Icon)
	level := strings.ToLower(strings.TrimSpace(req.Level))
	if level == "" {
		level = "easy"
	}
	if err := validateBadgeInput(name, level); err != nil {
		return admin.ManagementBadgeResponse{}, err
	}

	id, err := s.badgeRepo.Create(ctx, admin.CreateManagementBadgeRequest{
		Name:        name,
		Description: description,
		Icon:        icon,
		Level:       level,
	})
	if err != nil {
		return admin.ManagementBadgeResponse{}, err
	}
	return s.badgeRepo.GetByID(ctx, id, false)
}

func (s *managementBadgesService) UpdateBadge(ctx context.Context, id int64, req admin.UpdateManagementBadgeRequest) (admin.ManagementBadgeResponse, error) {
	if _, err := s.badgeRepo.GetByID(ctx, id, false); err != nil {
		return admin.ManagementBadgeResponse{}, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return admin.ManagementBadgeResponse{}, domain.NewBusinessError("INVALID_BADGE_NAME", "badge name cannot be empty", domain.ErrInvalidRequest)
		}
		req.Name = &name
	}
	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		req.Description = &description
	}
	if req.Icon != nil {
		icon := strings.TrimSpace(*req.Icon)
		req.Icon = &icon
	}
	if req.Level != nil {
		level := strings.ToLower(strings.TrimSpace(*req.Level))
		if !isValidBadgeLevel(level) {
			return admin.ManagementBadgeResponse{}, invalidBadgeLevelError()
		}
		req.Level = &level
	}

	if err := s.badgeRepo.Update(ctx, id, req); err != nil {
		return admin.ManagementBadgeResponse{}, err
	}
	return s.badgeRepo.GetByID(ctx, id, false)
}

func (s *managementBadgesService) SoftDeleteBadge(ctx context.Context, id int64) error {
	return s.badgeRepo.SoftDelete(ctx, id)
}

func (s *managementBadgesService) HardDeleteBadge(ctx context.Context, id int64) error {
	return s.badgeRepo.HardDelete(ctx, id)
}

func (s *managementBadgesService) RestoreBadge(ctx context.Context, id int64) (admin.ManagementBadgeResponse, error) {
	if err := s.badgeRepo.Restore(ctx, id); err != nil {
		return admin.ManagementBadgeResponse{}, err
	}
	return s.badgeRepo.GetByID(ctx, id, false)
}

func validateBadgeInput(name, level string) error {
	if strings.TrimSpace(name) == "" {
		return domain.NewBusinessError("INVALID_BADGE_NAME", "badge name cannot be empty", domain.ErrInvalidRequest)
	}
	if !isValidBadgeLevel(level) {
		return invalidBadgeLevelError()
	}
	return nil
}

func isValidBadgeLevel(level string) bool {
	switch level {
	case "easy", "medium", "hard":
		return true
	default:
		return false
	}
}

func invalidBadgeLevelError() error {
	return domain.NewBusinessError("INVALID_BADGE_LEVEL", "badge level must be easy, medium, or hard", domain.ErrInvalidRequest)
}
