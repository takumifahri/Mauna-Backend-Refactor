package soal

import (
	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	adminhandler.BaseHandler
	soalService usecase.ManagementSoalUsecase
}

func NewSoalHandler(soalService usecase.ManagementSoalUsecase) *Handler {
	return &Handler{
		BaseHandler: adminhandler.NewBaseHandler(),
		soalService: soalService,
	}
}
