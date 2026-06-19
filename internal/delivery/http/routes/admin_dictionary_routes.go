package routes

import (
	"net/http"
	"time"

	admindictionary "REFACTORING_MAUNA/internal/delivery/http/handler/admin/dictionary"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterAdminDictionaryRoutes(mux *http.ServeMux, dictService usecase.ManagementDictionaryUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	handler := admindictionary.NewDictionaryHandler(dictService)
	policy := usecase.RateLimitPolicy{Limit: 60, Window: time.Minute}
	protected := []middleware.Middleware{
		rateLimitMiddleware.Limiter(policy),
		middleware.Auth(tokenManager),
	}

	mux.Handle("GET /api/admin/dictionary", middleware.Chain(http.HandlerFunc(handler.ListDictionary), protected...))
	mux.Handle("POST /api/admin/dictionary", middleware.Chain(http.HandlerFunc(handler.CreateDictionary), protected...))
	mux.Handle("GET /api/admin/dictionary/{id}", middleware.Chain(http.HandlerFunc(handler.GetDictionary), protected...))
	mux.Handle("PATCH /api/admin/dictionary/{id}", middleware.Chain(http.HandlerFunc(handler.UpdateDictionary), protected...))
	mux.Handle("PATCH /api/admin/dictionary/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteDictionary), protected...))
	mux.Handle("DELETE /api/admin/dictionary/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteDictionary), protected...))
	mux.Handle("PATCH /api/admin/dictionary/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreDictionary), protected...))
}
