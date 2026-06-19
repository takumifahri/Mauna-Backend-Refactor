package usecase

import (
	"context"

	admin "REFACTORING_MAUNA/internal/dto/admin"
)

type ManagementDictionaryUsecase interface {
	ListDictionary(ctx context.Context, filter admin.ManagementDictionaryFilter) (admin.ManagementDictionaryListResponse, error)
	GetDictionary(ctx context.Context, id int64) (admin.ManagementDictionaryResponse, error)
	CreateDictionary(ctx context.Context, req admin.CreateManagementDictionaryRequest) (admin.ManagementDictionaryResponse, error)
	UpdateDictionary(ctx context.Context, id int64, req admin.UpdateManagementDictionaryRequest) (admin.ManagementDictionaryResponse, error)
	SoftDeleteDictionary(ctx context.Context, id int64) error
	HardDeleteDictionary(ctx context.Context, id int64) error
	RestoreDictionary(ctx context.Context, id int64) (admin.ManagementDictionaryResponse, error)
}
