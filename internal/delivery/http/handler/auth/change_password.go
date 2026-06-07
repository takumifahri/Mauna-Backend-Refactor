package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/pkg/validation"
)

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		logRequestError(r, "change_password_missing_auth_context", http.StatusUnauthorized, nil)
		writeMessageErrorResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logRequestError(r, "change_password_decode_failed", http.StatusBadRequest, err)
		writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}
	if err := validation.Validate(req); err != nil {
		logRequestError(r, "change_password_validation_failed", http.StatusBadRequest, err)
		writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.authService.ChangePassword(ctx, claims.UserID, req)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		logRequestError(r, "change_password_failed", statusCode, err, slog.String("user_id", claims.UserID))
		writeErrorResponse(w, statusCode, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Password changed successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
