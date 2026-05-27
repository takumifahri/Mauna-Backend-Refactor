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

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeProfileErrorResponse(w, http.StatusUnauthorized, domain.ErrUnauthorized)
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "profile_update_decode_failed", slog.Any("error", err))
		writeProfileErrorResponse(w, http.StatusBadRequest, domain.ErrInvalidRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.profileService.UpdateProfile(ctx, claims.UserID, req)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		writeProfileErrorResponse(w, statusCode, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Profile updated successfully",
		Data:      resp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
