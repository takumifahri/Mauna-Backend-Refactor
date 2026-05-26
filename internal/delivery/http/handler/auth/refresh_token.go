package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RefreshTokenRequest
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logRequestError(r, "refresh_token_decode_failed", http.StatusBadRequest, err)
			writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	refreshToken := getRefreshToken(r, req.RefreshToken)
	resp, err := h.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		logRequestError(r, "refresh_token_failed", statusCode, err)
		writeErrorResponse(w, statusCode, err)
		return
	}

	setAuthCookies(w, resp.AccessToken, resp.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Token refreshed successfully",
		Data:      resp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
