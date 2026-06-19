package routes

import (
	"net/http"
	"time"

	adminprogress "REFACTORING_MAUNA/internal/delivery/http/handler/admin/progress"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterAdminProgressRoutes(mux *http.ServeMux, progressService usecase.ManagementProgressUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	handler := adminprogress.NewProgressHandler(progressService)
	policy := usecase.RateLimitPolicy{Limit: 60, Window: time.Minute}
	protected := []middleware.Middleware{
		rateLimitMiddleware.Limiter(policy),
		middleware.Auth(tokenManager),
	}

	mux.Handle("GET /api/admin/progress", middleware.Chain(http.HandlerFunc(handler.ListProgress), protected...))
	mux.Handle("GET /api/admin/progress/{id}", middleware.Chain(http.HandlerFunc(handler.GetProgress), protected...))
	mux.Handle("PATCH /api/admin/progress/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteProgress), protected...))
	mux.Handle("DELETE /api/admin/progress/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteProgress), protected...))
	mux.Handle("PATCH /api/admin/progress/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreProgress), protected...))
}
