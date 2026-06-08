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

	mux.Handle("GET /api/admin/users", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.ListUsers))))
	mux.Handle("POST /api/admin/users", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.CreateUser))))
	mux.Handle("GET /api/admin/users/{id}", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.GetUser))))
	mux.Handle("PATCH /api/admin/users/{id}", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.UpdateUser))))
	mux.Handle("PATCH /api/admin/users/{id}/soft-delete", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.SoftDeleteUser))))
	mux.Handle("DELETE /api/admin/users/{id}", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.DeleteUser))))
	mux.Handle("PATCH /api/admin/users/{id}/restore", rateLimitMiddleware.Limit(policy, middleware.JWTAuth(tokenManager, http.HandlerFunc(handler.RestoreUser))))
}
