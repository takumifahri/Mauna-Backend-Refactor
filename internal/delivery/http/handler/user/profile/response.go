package profile

import (
	"encoding/json"
	"net/http"

	"REFACTORING_MAUNA/internal/dto"
)

func writeProfileErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: err.Error(),
	})
}
