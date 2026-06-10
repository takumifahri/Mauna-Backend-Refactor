package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/delivery/http/middleware"
	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

type BaseHandler struct{}

func NewBaseHandler() BaseHandler {
	return BaseHandler{}
}

func (h BaseHandler) RequireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		h.WriteError(w, domain.ErrUnauthorized)
		return false
	}
	if claims.Role != "admin" {
		h.WriteError(w, domain.ErrForbidden)
		return false
	}
	return true
}

func (BaseHandler) WriteJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func (BaseHandler) WriteError(w http.ResponseWriter, err error) {
	statusCode := domain.ErrorToStatusCode(err)
	if statusCode == http.StatusOK {
		statusCode = http.StatusInternalServerError
	}

	var businessErr domain.BusinessError
	message := err.Error()
	if errors.As(err, &businessErr) {
		message = businessErr.Message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   domain.DebugError(err),
	})
}
