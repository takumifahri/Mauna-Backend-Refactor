package auth

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "REFACTORING_MAUNA/internal/domain"
    "REFACTORING_MAUNA/internal/dto"
)

// ← Handler method (HTTP layer)
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req dto.LoginRequest
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

    // ← Call service (business logic)
    resp, err := h.authService.Login(ctx, req)
    if err != nil {
        statusCode := domain.ErrorToStatusCode(err)
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
        Message:   "Login successful",
        Data:      resp,
        Timestamp: time.Now().Format(time.RFC3339),
    })
}