package usecase

import (
	"context"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

type ManagementLevelsUsecase interface {
	ListLevels(ctx context.Context, filter admin.ManagementLevelsFilter) (admin.ManagementLevelsListResponse, error)
	GetLevel(ctx context.Context, id int64) (admin.ManagementLevelResponse, error)
	CreateLevel(ctx context.Context, req admin.CreateManagementLevelRequest) (admin.ManagementLevelResponse, error)
	UpdateLevel(ctx context.Context, id int64, req admin.UpdateManagementLevelRequest) (admin.ManagementLevelResponse, error)
	SoftDeleteLevel(ctx context.Context, id int64) error
	HardDeleteLevel(ctx context.Context, id int64) error
	RestoreLevel(ctx context.Context, id int64) (admin.ManagementLevelResponse, error)
}

type ManagementSubLevelsUsecase interface {
	ListSubLevels(ctx context.Context, filter admin.ManagementSubLevelsFilter) (admin.ManagementSubLevelsListResponse, error)
	GetSubLevel(ctx context.Context, id int64) (admin.ManagementSubLevelResponse, error)
	CreateSubLevel(ctx context.Context, req admin.CreateManagementSubLevelRequest) (admin.ManagementSubLevelResponse, error)
	UpdateSubLevel(ctx context.Context, id int64, req admin.UpdateManagementSubLevelRequest) (admin.ManagementSubLevelResponse, error)
	SoftDeleteSubLevel(ctx context.Context, id int64) error
	HardDeleteSubLevel(ctx context.Context, id int64) error
	RestoreSubLevel(ctx context.Context, id int64) (admin.ManagementSubLevelResponse, error)
}
