package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/pkg/database"
)

// HealthHandler checks database and server health
func HealthHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := dto.Response{ // ← Pakai dari dto
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		// Check database
		if err := db.Health(); err != nil {
			resp.Status = "unhealthy"
			resp.Message = fmt.Sprintf("Database error: %v", err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp.Data = map[string]string{
			"database": "connected",
			"server":   "running",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// RootHandler menampilkan info API
func RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := dto.Response{ // ← Pakai dari dto
			Status:    "success",
			Message:   "Mauna Backend API",
			Timestamp: time.Now().Format(time.RFC3339),
			Data: map[string]string{
				"version":     "1.0.0",
				"name":        "Mauna Backend",
				"description": "Sign Language Learning API",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
