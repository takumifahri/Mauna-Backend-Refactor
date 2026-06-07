package usecase

import (
	"context"

	"REFACTORING_MAUNA/internal/dto/admin"
)

type ManagementUsersUsecase interface {
	ListUsers(ctx context.Context, filter admin.ManagementUsersFilter) (admin.ManagementUsersListResponse, error)
	GetUser(ctx context.Context, id string) (admin.ManagementUserResponse, error)
	CreateUser(ctx context.Context, req admin.CreateManagementUserRequest) (admin.ManagementUserResponse, error)
	UpdateUser(ctx context.Context, id string, req admin.UpdateManagementUserRequest) (admin.ManagementUserResponse, error)
	SoftDeleteUser(ctx context.Context, id string) error
	HardDeleteUser(ctx context.Context, id string) error
	RestoreUser(ctx context.Context, id string) (admin.ManagementUserResponse, error)
}
