package profile

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeProfileErrorResponse(w, http.StatusUnauthorized, domain.ErrUnauthorized)
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "profile_change_password_decode_failed", slog.Any("error", err))
		writeProfileErrorResponse(w, http.StatusBadRequest, domain.ErrInvalidRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.profileService.ChangePassword(ctx, claims.UserID, req); err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		writeProfileErrorResponse(w, statusCode, err)
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
