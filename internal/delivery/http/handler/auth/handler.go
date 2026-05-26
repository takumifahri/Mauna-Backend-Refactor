package auth

import (
    "REFACTORING_MAUNA/internal/usecase"
)

// ← Handler struct untuk HTTP
type Handler struct {
    authService usecase.AuthUsecase  // Inject service interface
}

func NewAuthHandler(authService usecase.AuthUsecase) *Handler {
    return &Handler{
        authService: authService,
    }
}