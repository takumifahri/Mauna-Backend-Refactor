package routes

import (
    "net/http"

    "REFACTORING_MAUNA/internal/delivery/http/handler/auth"
    "REFACTORING_MAUNA/internal/usecase"
    // "REFACTORING_MAUNA/pkg/database"
)

// RegisterAuthRoutes registers all auth-related routes
func RegisterAuthRoutes(mux *http.ServeMux, authService usecase.AuthUsecase) {
    authHandler := auth.NewAuthHandler(authService)

    // Auth endpoints
    mux.HandleFunc("POST /api/auth/login", authHandler.Login)
    mux.HandleFunc("POST /api/auth/register", authHandler.Register)
    mux.HandleFunc("POST /api/auth/change-password", authHandler.ChangePassword)
    mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
    mux.HandleFunc("POST /api/auth/refresh-token", authHandler.RefreshToken)
}