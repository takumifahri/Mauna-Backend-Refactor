package progress

import (
	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	adminhandler.BaseHandler
	progressService usecase.ManagementProgressUsecase
}

func NewProgressHandler(progressService usecase.ManagementProgressUsecase) *Handler {
	return &Handler{
		BaseHandler:     adminhandler.NewBaseHandler(),
		progressService: progressService,
	}
}
