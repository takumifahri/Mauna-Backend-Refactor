package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

func writeAdminJSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.Response{
		Status:    "success",
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func writeAdminError(w http.ResponseWriter, err error) {
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
