package routes

import (
	"net/http"
	"time"

	adminusers "REFACTORING_MAUNA/internal/delivery/http/handler/admin/users"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

func RegisterAdminUsersRoutes(mux *http.ServeMux, userService usecase.ManagementUsersUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	handler := adminusers.NewUserHandler(userService)
	policy := usecase.RateLimitPolicy{Limit: 60, Window: time.Minute}
	protected := []middleware.Middleware{
		rateLimitMiddleware.Limiter(policy),
		middleware.Auth(tokenManager),
	}

	mux.Handle("GET /api/admin/users", middleware.Chain(http.HandlerFunc(handler.ListUsers), protected...))
	mux.Handle("POST /api/admin/users", middleware.Chain(http.HandlerFunc(handler.CreateUser), protected...))
	mux.Handle("GET /api/admin/users/{id}", middleware.Chain(http.HandlerFunc(handler.GetUser), protected...))
	mux.Handle("PATCH /api/admin/users/{id}", middleware.Chain(http.HandlerFunc(handler.UpdateUser), protected...))
	mux.Handle("PATCH /api/admin/users/{id}/soft-delete", middleware.Chain(http.HandlerFunc(handler.SoftDeleteUser), protected...))
	mux.Handle("DELETE /api/admin/users/{id}", middleware.Chain(http.HandlerFunc(handler.DeleteUser), protected...))
	mux.Handle("PATCH /api/admin/users/{id}/restore", middleware.Chain(http.HandlerFunc(handler.RestoreUser), protected...))
}
