package usecase

import (
	"context"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

type ManagementProgressUsecase interface {
	ListProgress(ctx context.Context, filter admin.ManagementProgressFilter) (admin.ManagementProgressListResponse, error)
	GetProgress(ctx context.Context, id int64) (admin.ManagementProgressResponse, error)
	SoftDeleteProgress(ctx context.Context, id int64) error
	HardDeleteProgress(ctx context.Context, id int64) error
	RestoreProgress(ctx context.Context, id int64) (admin.ManagementProgressResponse, error)
}
