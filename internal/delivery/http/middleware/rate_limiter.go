package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/internal/usecase"
)

type RateLimitMiddleware struct {
	rateLimiter usecase.RateLimiterUsecase
}

func NewRateLimitMiddleware(rateLimiter usecase.RateLimiterUsecase) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rateLimiter: rateLimiter,
	}
}

func (m *RateLimitMiddleware) Limit(policy usecase.RateLimitPolicy, next http.Handler) http.Handler {
	return m.Limiter(policy)(next)
}

func (m *RateLimitMiddleware) Limiter(policy usecase.RateLimitPolicy) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			decision, err := m.rateLimiter.Allow(r.Context(), rateLimitKey(r), policy)
			if err != nil {
				logRequestError(r, "rate_limit_failed", http.StatusInternalServerError, err)
				writeRateLimitError(w)
				return
			}
			if decision.Allowed {
				next.ServeHTTP(w, r)
				return
			}

			logRequestError(r, "rate_limit_exceeded", http.StatusTooManyRequests, nil)
			writeRateLimitExceeded(w, decision.RetryAfter)
		})
	}
}

func rateLimitKey(r *http.Request) string {
	path := r.Pattern
	if path == "" {
		path = r.URL.Path
	}
	if strings.HasPrefix(path, r.Method+" ") {
		path = strings.TrimPrefix(path, r.Method+" ")
	}

	return r.Method + ":" + path + ":" + clientIP(r)
}

func clientIP(r *http.Request) string {
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func writeRateLimitExceeded(w http.ResponseWriter, retryAfter time.Duration) {
	retryAfterSeconds := int(retryAfter.Seconds())
	if retryAfterSeconds < 1 {
		retryAfterSeconds = 1
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", strconv.Itoa(retryAfterSeconds))
	w.WriteHeader(http.StatusTooManyRequests)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: "Too many requests, please try again later",
	})
}

func writeRateLimitError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: "internal server error",
	})
}
