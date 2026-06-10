package usecase

import (
	"context"

	"REFACTORING_MAUNA/internal/dto/admin"
)

type ManagementBadgesUsecase interface {
	ListBadges(ctx context.Context, filter admin.ManagementBadgesFilter) (admin.ManagementBadgesListResponse, error)
	GetBadge(ctx context.Context, id int64) (admin.ManagementBadgeResponse, error)
	CreateBadge(ctx context.Context, req admin.CreateManagementBadgeRequest) (admin.ManagementBadgeResponse, error)
	UpdateBadge(ctx context.Context, id int64, req admin.UpdateManagementBadgeRequest) (admin.ManagementBadgeResponse, error)
	SoftDeleteBadge(ctx context.Context, id int64) error
	HardDeleteBadge(ctx context.Context, id int64) error
	RestoreBadge(ctx context.Context, id int64) (admin.ManagementBadgeResponse, error)
}
