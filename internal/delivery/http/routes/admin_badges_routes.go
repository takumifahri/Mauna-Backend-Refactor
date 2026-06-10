package routes

import (
	"net/http"
	"time"

	adminbadges "REFACTORING_MAUNA/internal/delivery/http/handler/admin/badges"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterAdminBadgesRoutes(mux *http.ServeMux, badgeService usecase.ManagementBadgesUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	handler := adminbadges.NewBadgeHandler(badgeService)
	policy := usecase.RateLimitPolicy{Limit: 60, Window: time.Minute}
	protected := []middleware.Middleware{
		rateLimitMiddleware.Limiter(policy),
		middleware.Auth(tokenManager),
	}

	mux.Handle("GET /api/admin/badges", middleware.Chain(http.HandlerFunc(handler.ListBadges), protected...))
	mux.Handle("POST /api/admin/badges", middleware.Chain(http.HandlerFunc(handler.CreateBadge), protected...))
	mux.Handle("GET /api/admin/badges/{id}", middleware.Chain(http.HandlerFunc(handler.GetBadge), protected...))
	mux.Handle("PATCH /api/admin/badges/{id}", middleware.Chain(http.HandlerFunc(handler.UpdateBadge), protected...))
	mux.Handle("PATCH /api/admin/badges/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteBadge), protected...))
	mux.Handle("DELETE /api/admin/badges/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteBadge), protected...))
	mux.Handle("PATCH /api/admin/badges/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreBadge), protected...))
}
