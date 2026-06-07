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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logRequestError(r, "register_decode_failed", http.StatusBadRequest, err)
		writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}
	if err := validation.Validate(req); err != nil {
		logRequestError(r, "register_validation_failed", http.StatusBadRequest, err)
		writeMessageErrorResponse(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.authService.Register(ctx, req)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		logRequestError(r, "register_failed", statusCode, err)
		writeErrorResponse(w, statusCode, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   "Registration successful",
		Data:      resp,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
