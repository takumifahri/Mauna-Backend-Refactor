package users

import (
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	userService usecase.ManagementUsersUsecase
}

func NewUserHandler(userService usecase.ManagementUsersUsecase) *Handler {
	return &Handler{userService: userService}
}
