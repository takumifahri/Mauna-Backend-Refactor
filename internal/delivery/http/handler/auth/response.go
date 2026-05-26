package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/dto"
)

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: err.Error(),
		Error:   debugError(err),
	})
}

func writeMessageErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   debugError(err),
	})
}

func debugError(err error) string {
	if !appDebugEnabled() {
		return ""
	}
	return domain.DebugError(err)
}

func appDebugEnabled() bool {
	value := strings.ToLower(os.Getenv("APP_DEBUG"))
	return value == "true" || value == "1" || value == "yes"
}
