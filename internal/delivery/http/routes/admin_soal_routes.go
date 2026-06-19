package routes

import (
	"net/http"
	"time"

	adminsoal "REFACTORING_MAUNA/internal/delivery/http/handler/admin/soal"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterAdminSoalRoutes(mux *http.ServeMux, soalService usecase.ManagementSoalUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	handler := adminsoal.NewSoalHandler(soalService)
	policy := usecase.RateLimitPolicy{Limit: 60, Window: time.Minute}
	protected := []middleware.Middleware{
		rateLimitMiddleware.Limiter(policy),
		middleware.Auth(tokenManager),
	}

	mux.Handle("GET /api/admin/soal", middleware.Chain(http.HandlerFunc(handler.ListSoal), protected...))
	mux.Handle("POST /api/admin/soal", middleware.Chain(http.HandlerFunc(handler.CreateSoal), protected...))
	mux.Handle("GET /api/admin/soal/{id}", middleware.Chain(http.HandlerFunc(handler.GetSoal), protected...))
	mux.Handle("PATCH /api/admin/soal/{id}", middleware.Chain(http.HandlerFunc(handler.UpdateSoal), protected...))
	mux.Handle("PATCH /api/admin/soal/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteSoal), protected...))
	mux.Handle("DELETE /api/admin/soal/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteSoal), protected...))
	mux.Handle("PATCH /api/admin/soal/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreSoal), protected...))
}
