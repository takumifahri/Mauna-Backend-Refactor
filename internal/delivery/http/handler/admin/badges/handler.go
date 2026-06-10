package badges

import (
	adminhandler "REFACTORING_MAUNA/internal/delivery/http/handler/admin"
	"REFACTORING_MAUNA/internal/usecase"
)

type Handler struct {
	adminhandler.BaseHandler
	badgeService usecase.ManagementBadgesUsecase
}

func NewBadgeHandler(badgeService usecase.ManagementBadgesUsecase) *Handler {
	return &Handler{
		BaseHandler:  adminhandler.NewBaseHandler(),
		badgeService: badgeService,
	}
}
