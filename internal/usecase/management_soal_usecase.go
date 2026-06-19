package usecase

import (
	"context"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

type ManagementSoalUsecase interface {
	ListSoal(ctx context.Context, filter admin.ManagementSoalFilter) (admin.ManagementSoalListResponse, error)
	GetSoal(ctx context.Context, id int64) (admin.ManagementSoalResponse, error)
	CreateSoal(ctx context.Context, req admin.CreateManagementSoalRequest) (admin.ManagementSoalResponse, error)
	UpdateSoal(ctx context.Context, id int64, req admin.UpdateManagementSoalRequest) (admin.ManagementSoalResponse, error)
	SoftDeleteSoal(ctx context.Context, id int64) error
	HardDeleteSoal(ctx context.Context, id int64) error
	RestoreSoal(ctx context.Context, id int64) (admin.ManagementSoalResponse, error)
}
