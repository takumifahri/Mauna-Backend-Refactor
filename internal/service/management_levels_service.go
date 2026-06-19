package service

import (
	"context"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

// ─── Level ────────────────────────────────────────────────────────────────────

type managementLevelsService struct {
	levelRepo domain.ManagementLevelsRepository
}

var _ usecase.ManagementLevelsUsecase = (*managementLevelsService)(nil)

func NewManagementLevelsService(levelRepo domain.ManagementLevelsRepository) usecase.ManagementLevelsUsecase {
	return &managementLevelsService{levelRepo: levelRepo}
}

func (s *managementLevelsService) ListLevels(ctx context.Context, filter admin.ManagementLevelsFilter) (admin.ManagementLevelsListResponse, error) {
	return s.levelRepo.List(ctx, filter)
}

func (s *managementLevelsService) GetLevel(ctx context.Context, id int64) (admin.ManagementLevelResponse, error) {
	return s.levelRepo.GetByID(ctx, id, false)
}

func (s *managementLevelsService) CreateLevel(ctx context.Context, req admin.CreateManagementLevelRequest) (admin.ManagementLevelResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return admin.ManagementLevelResponse{}, domain.NewBusinessError("INVALID_LEVEL_NAME", "level name cannot be empty", domain.ErrInvalidRequest)
	}

	id, err := s.levelRepo.Create(ctx, admin.CreateManagementLevelRequest{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Tujuan:      strings.TrimSpace(req.Tujuan),
	})
	if err != nil {
		return admin.ManagementLevelResponse{}, err
	}
	return s.levelRepo.GetByID(ctx, id, false)
}

func (s *managementLevelsService) UpdateLevel(ctx context.Context, id int64, req admin.UpdateManagementLevelRequest) (admin.ManagementLevelResponse, error) {
	if _, err := s.levelRepo.GetByID(ctx, id, false); err != nil {
		return admin.ManagementLevelResponse{}, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return admin.ManagementLevelResponse{}, domain.NewBusinessError("INVALID_LEVEL_NAME", "level name cannot be empty", domain.ErrInvalidRequest)
		}
		req.Name = &name
	}
	if req.Description != nil {
		v := strings.TrimSpace(*req.Description)
		req.Description = &v
	}
	if req.Tujuan != nil {
		v := strings.TrimSpace(*req.Tujuan)
		req.Tujuan = &v
	}

	if err := s.levelRepo.Update(ctx, id, req); err != nil {
		return admin.ManagementLevelResponse{}, err
	}
	return s.levelRepo.GetByID(ctx, id, false)
}

func (s *managementLevelsService) SoftDeleteLevel(ctx context.Context, id int64) error {
	return s.levelRepo.SoftDelete(ctx, id)
}

func (s *managementLevelsService) HardDeleteLevel(ctx context.Context, id int64) error {
	return s.levelRepo.HardDelete(ctx, id)
}

func (s *managementLevelsService) RestoreLevel(ctx context.Context, id int64) (admin.ManagementLevelResponse, error) {
	if err := s.levelRepo.Restore(ctx, id); err != nil {
		return admin.ManagementLevelResponse{}, err
	}
	return s.levelRepo.GetByID(ctx, id, false)
}

// ─── SubLevel ─────────────────────────────────────────────────────────────────

type managementSubLevelsService struct {
	subLevelRepo domain.ManagementSubLevelsRepository
}

var _ usecase.ManagementSubLevelsUsecase = (*managementSubLevelsService)(nil)

func NewManagementSubLevelsService(subLevelRepo domain.ManagementSubLevelsRepository) usecase.ManagementSubLevelsUsecase {
	return &managementSubLevelsService{subLevelRepo: subLevelRepo}
}

func (s *managementSubLevelsService) ListSubLevels(ctx context.Context, filter admin.ManagementSubLevelsFilter) (admin.ManagementSubLevelsListResponse, error) {
	return s.subLevelRepo.List(ctx, filter)
}

func (s *managementSubLevelsService) GetSubLevel(ctx context.Context, id int64) (admin.ManagementSubLevelResponse, error) {
	return s.subLevelRepo.GetByID(ctx, id, false)
}

func (s *managementSubLevelsService) CreateSubLevel(ctx context.Context, req admin.CreateManagementSubLevelRequest) (admin.ManagementSubLevelResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return admin.ManagementSubLevelResponse{}, domain.NewBusinessError("INVALID_SUBLEVEL_NAME", "sublevel name cannot be empty", domain.ErrInvalidRequest)
	}
	if req.LevelID <= 0 {
		return admin.ManagementSubLevelResponse{}, domain.NewBusinessError("INVALID_LEVEL_ID", "level_id is required", domain.ErrInvalidRequest)
	}

	id, err := s.subLevelRepo.Create(ctx, admin.CreateManagementSubLevelRequest{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		Tujuan:      strings.TrimSpace(req.Tujuan),
		LevelID:     req.LevelID,
	})
	if err != nil {
		return admin.ManagementSubLevelResponse{}, err
	}
	return s.subLevelRepo.GetByID(ctx, id, false)
}

func (s *managementSubLevelsService) UpdateSubLevel(ctx context.Context, id int64, req admin.UpdateManagementSubLevelRequest) (admin.ManagementSubLevelResponse, error) {
	if _, err := s.subLevelRepo.GetByID(ctx, id, false); err != nil {
		return admin.ManagementSubLevelResponse{}, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return admin.ManagementSubLevelResponse{}, domain.NewBusinessError("INVALID_SUBLEVEL_NAME", "sublevel name cannot be empty", domain.ErrInvalidRequest)
		}
		req.Name = &name
	}
	if req.Description != nil {
		v := strings.TrimSpace(*req.Description)
		req.Description = &v
	}
	if req.Tujuan != nil {
		v := strings.TrimSpace(*req.Tujuan)
		req.Tujuan = &v
	}
	if req.LevelID != nil && *req.LevelID <= 0 {
		return admin.ManagementSubLevelResponse{}, domain.NewBusinessError("INVALID_LEVEL_ID", "level_id must be positive", domain.ErrInvalidRequest)
	}

	if err := s.subLevelRepo.Update(ctx, id, req); err != nil {
		return admin.ManagementSubLevelResponse{}, err
	}
	return s.subLevelRepo.GetByID(ctx, id, false)
}

func (s *managementSubLevelsService) SoftDeleteSubLevel(ctx context.Context, id int64) error {
	return s.subLevelRepo.SoftDelete(ctx, id)
}

func (s *managementSubLevelsService) HardDeleteSubLevel(ctx context.Context, id int64) error {
	return s.subLevelRepo.HardDelete(ctx, id)
}

func (s *managementSubLevelsService) RestoreSubLevel(ctx context.Context, id int64) (admin.ManagementSubLevelResponse, error) {
	if err := s.subLevelRepo.Restore(ctx, id); err != nil {
		return admin.ManagementSubLevelResponse{}, err
	}
	return s.subLevelRepo.GetByID(ctx, id, false)
}
