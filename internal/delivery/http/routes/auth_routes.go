package routes

import (
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/handler/auth"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
)

// RegisterAuthRoutes registers all auth-related routes
func RegisterAuthRoutes(mux *http.ServeMux, authService usecase.AuthUsecase, tokenManager usecase.TokenManager, rateLimitMiddleware *middleware.RateLimitMiddleware) {
	authHandler := auth.NewAuthHandler(authService)

	// Auth endpoints
	mux.Handle("POST /api/auth/login", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute}, http.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /api/auth/register", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 3, Window: time.Minute}, http.HandlerFunc(authHandler.Register)))
	mux.Handle("POST /api/auth/change-password", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute}, middleware.JWTAuth(tokenManager, http.HandlerFunc(authHandler.ChangePassword))))
	mux.Handle("POST /api/auth/logout", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 20, Window: time.Minute}, http.HandlerFunc(authHandler.Logout)))
	mux.Handle("POST /api/auth/refresh-token", rateLimitMiddleware.Limit(usecase.RateLimitPolicy{Limit: 10, Window: time.Minute}, http.HandlerFunc(authHandler.RefreshToken)))
}
