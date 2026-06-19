package dictionary

import (
	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	adminhandler.BaseHandler
	dictService usecase.ManagementDictionaryUsecase
}

func NewDictionaryHandler(dictService usecase.ManagementDictionaryUsecase) *Handler {
	return &Handler{
		BaseHandler: adminhandler.NewBaseHandler(),
		dictService: dictService,
	}
}
