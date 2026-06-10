package users

import (
	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	adminhandler.BaseHandler
	userService usecase.ManagementUsersUsecase
}

func NewUserHandler(userService usecase.ManagementUsersUsecase) *Handler {
	return &Handler{
		BaseHandler: adminhandler.NewBaseHandler(),
		userService: userService,
	}
}
