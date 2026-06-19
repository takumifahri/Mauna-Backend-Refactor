package levels

import (
	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	adminhandler.BaseHandler
	levelService    usecase.ManagementLevelsUsecase
	subLevelService usecase.ManagementSubLevelsUsecase
}

func NewLevelHandler(levelService usecase.ManagementLevelsUsecase, subLevelService usecase.ManagementSubLevelsUsecase) *Handler {
	return &Handler{
		BaseHandler:     adminhandler.NewBaseHandler(),
		levelService:    levelService,
		subLevelService: subLevelService,
	}
}
