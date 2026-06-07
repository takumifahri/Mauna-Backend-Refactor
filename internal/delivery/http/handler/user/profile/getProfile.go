package profile

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.profileService.GetProfile(ctx, claims.UserID)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Profile retrieved successfully",
		Data:      resp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
