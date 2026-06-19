package service

import (
	"context"

	"REFACTORING_MAUNA/internal/domain"
	admin "REFACTORING_MAUNA/internal/dto/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type managementProgressService struct {
	progressRepo domain.ManagementProgressRepository
}

var _ usecase.ManagementProgressUsecase = (*managementProgressService)(nil)

func NewManagementProgressService(progressRepo domain.ManagementProgressRepository) usecase.ManagementProgressUsecase {
	return &managementProgressService{progressRepo: progressRepo}
}

func (s *managementProgressService) ListProgress(ctx context.Context, filter admin.ManagementProgressFilter) (admin.ManagementProgressListResponse, error) {
	return s.progressRepo.List(ctx, filter)
}

func (s *managementProgressService) GetProgress(ctx context.Context, id int64) (admin.ManagementProgressResponse, error) {
	return s.progressRepo.GetByID(ctx, id, false)
}

func (s *managementProgressService) SoftDeleteProgress(ctx context.Context, id int64) error {
	return s.progressRepo.SoftDelete(ctx, id)
}

func (s *managementProgressService) HardDeleteProgress(ctx context.Context, id int64) error {
	return s.progressRepo.HardDelete(ctx, id)
}

func (s *managementProgressService) RestoreProgress(ctx context.Context, id int64) (admin.ManagementProgressResponse, error) {
	if err := s.progressRepo.Restore(ctx, id); err != nil {
		return admin.ManagementProgressResponse{}, err
	}
	return s.progressRepo.GetByID(ctx, id, false)
}
