package auth

import (
    "REFACTORING_MAUNA/internal/service"
)

// ← Handler struct untuk HTTP
type Handler struct {
    authService service.AuthUsecase  // Inject service interface
}

func NewAuthHandler(authService service.AuthUsecase) *Handler {
    return &Handler{
        authService: authService,
    }
}