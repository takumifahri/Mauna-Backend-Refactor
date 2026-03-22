package auth

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "REFACTORING_MAUNA/internal/domain"
    "REFACTORING_MAUNA/internal/dto"
)

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get userID from JWT (middleware nanti)
    userIDStr := r.Header.Get("X-User-ID")
    if userIDStr == "" {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(dto.ErrorResponse{
            Status:  "error",
            Message: "Unauthorized",
        })
        return
    }

    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(dto.ErrorResponse{
            Status:  "error",
            Message: "Invalid user ID",
        })
        return
    }

    var req dto.ChangePasswordRequest
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

    err = h.authService.ChangePassword(ctx, userID, req)
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
        Message:   "Password changed successfully",
        Timestamp: time.Now().Format(time.RFC3339),
    })
}