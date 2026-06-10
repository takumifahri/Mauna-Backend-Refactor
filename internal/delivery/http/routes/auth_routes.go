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
	authn := middleware.Auth(tokenManager)
	limit := rateLimitMiddleware.Limiter

	// Auth endpoints
	mux.Handle("POST /api/auth/login", middleware.Chain(http.HandlerFunc(authHandler.Login), limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute})))
	mux.Handle("POST /api/auth/register", middleware.Chain(http.HandlerFunc(authHandler.Register), limit(usecase.RateLimitPolicy{Limit: 3, Window: time.Minute})))
	mux.Handle("POST /api/auth/verify-registration", middleware.Chain(http.HandlerFunc(authHandler.VerifyRegistration), limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute})))
	mux.Handle("POST /api/auth/forgot-password", middleware.Chain(http.HandlerFunc(authHandler.ForgotPassword), limit(usecase.RateLimitPolicy{Limit: 3, Window: 15 * time.Minute})))
	mux.Handle("POST /api/auth/reset-password", middleware.Chain(http.HandlerFunc(authHandler.ResetPassword), limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute})))
	mux.Handle("POST /api/auth/change-password", middleware.Chain(http.HandlerFunc(authHandler.ChangePassword), limit(usecase.RateLimitPolicy{Limit: 5, Window: time.Minute}), authn))
	mux.Handle("POST /api/auth/logout", middleware.Chain(http.HandlerFunc(authHandler.Logout), limit(usecase.RateLimitPolicy{Limit: 20, Window: time.Minute})))
	mux.Handle("POST /api/auth/refresh-token", middleware.Chain(http.HandlerFunc(authHandler.RefreshToken), limit(usecase.RateLimitPolicy{Limit: 10, Window: time.Minute})))
}
