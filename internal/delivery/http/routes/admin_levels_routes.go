package routes

import (
	"net/http"
	"time"

	adminlevels "REFACTORING_MAUNA/internal/delivery/http/handler/admin/levels"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterAdminLevelsRoutes(mux *http.ServeMux, levelService usecase.ManagementLevelsUsecase, subLevelService usecase.ManagementSubLevelsUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	handler := adminlevels.NewLevelHandler(levelService, subLevelService)
	policy := usecase.RateLimitPolicy{Limit: 60, Window: time.Minute}
	protected := []middleware.Middleware{
		rateLimitMiddleware.Limiter(policy),
		middleware.Auth(tokenManager),
	}

	// Level routes
	mux.Handle("GET /api/admin/levels", middleware.Chain(http.HandlerFunc(handler.ListLevels), protected...))
	mux.Handle("POST /api/admin/levels", middleware.Chain(http.HandlerFunc(handler.CreateLevel), protected...))
	mux.Handle("GET /api/admin/levels/{id}", middleware.Chain(http.HandlerFunc(handler.GetLevel), protected...))
	mux.Handle("PATCH /api/admin/levels/{id}", middleware.Chain(http.HandlerFunc(handler.UpdateLevel), protected...))
	mux.Handle("PATCH /api/admin/levels/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteLevel), protected...))
	mux.Handle("DELETE /api/admin/levels/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteLevel), protected...))
	mux.Handle("PATCH /api/admin/levels/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreLevel), protected...))

	// SubLevel routes
	mux.Handle("GET /api/admin/sublevels", middleware.Chain(http.HandlerFunc(handler.ListSubLevels), protected...))
	mux.Handle("POST /api/admin/sublevels", middleware.Chain(http.HandlerFunc(handler.CreateSubLevel), protected...))
	mux.Handle("GET /api/admin/sublevels/{id}", middleware.Chain(http.HandlerFunc(handler.GetSubLevel), protected...))
	mux.Handle("PATCH /api/admin/sublevels/{id}", middleware.Chain(http.HandlerFunc(handler.UpdateSubLevel), protected...))
	mux.Handle("PATCH /api/admin/sublevels/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteSubLevel), protected...))
	mux.Handle("DELETE /api/admin/sublevels/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteSubLevel), protected...))
	mux.Handle("PATCH /api/admin/sublevels/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreSubLevel), protected...))
}
