package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LogoutRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logRequestError(r, "logout_decode_failed", http.StatusBadRequest, err)
			writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	refreshToken := getRefreshToken(r, req.RefreshToken)
	err := h.authService.Logout(ctx, refreshToken)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		logRequestError(r, "logout_failed", statusCode, err)
		writeErrorResponse(w, statusCode, err)
		return
	}

	clearAuthCookies(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	})
}
