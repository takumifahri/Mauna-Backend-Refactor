package routes

import (
	"net/http"

	"REFACTORING_MAUNA/internal/delivery/http/handler/auth"
	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/usecase"
	// "REFACTORING_MAUNA/pkg/database"
)

// RegisterAuthRoutes registers all auth-related routes
func RegisterAuthRoutes(mux *http.ServeMux, authService usecase.AuthUsecase, tokenManager usecase.TokenManager) {
	authHandler := auth.NewAuthHandler(authService)

	// Auth endpoints
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.Handle("POST /api/auth/change-password", middleware.JWTAuth(tokenManager, http.HandlerFunc(authHandler.ChangePassword)))
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/auth/refresh-token", authHandler.RefreshToken)
}
