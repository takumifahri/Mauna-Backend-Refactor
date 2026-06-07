package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/pkg/validation"
)

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logRequestError(r, "reset_password_decode_failed", http.StatusBadRequest, err)
		writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}
	if err := validation.Validate(req); err != nil {
		logRequestError(r, "reset_password_validation_failed", http.StatusBadRequest, err)
		writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.authService.ResetPassword(ctx, req); err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		logRequestError(r, "reset_password_failed", statusCode, err)
		writeErrorResponse(w, statusCode, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Password reset successfully",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
