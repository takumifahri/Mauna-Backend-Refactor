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

	var req dto.LogoutResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Status:  "error",
			Message: "Invalid request",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		statusCode := domain.ErrorToStatusCode(err)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
}