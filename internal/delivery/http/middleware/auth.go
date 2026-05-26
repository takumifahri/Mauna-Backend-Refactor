package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"REFACTORING_MAUNA/internal/dto"
	"REFACTORING_MAUNA/internal/usecase"
)

type authContextKey string

const (
	authClaimsKey         authContextKey = "auth_claims"
	accessTokenCookieName string         = "access_token"
)

func JWTAuth(tokenManager usecase.TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := accessTokenFromRequest(r)
		if token == "" {
			logRequestError(r, "auth_token_missing", http.StatusUnauthorized, nil)
			writeUnauthorized(w)
			return
		}

		claims, err := tokenManager.VerifyToken(token)
		if err != nil {
			logRequestError(r, "auth_token_invalid", http.StatusUnauthorized, err, slog.String("token_source", tokenSource(r)))
			writeUnauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), authClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromContext(ctx context.Context) (usecase.TokenClaims, bool) {
	claims, ok := ctx.Value(authClaimsKey).(usecase.TokenClaims)
	return claims, ok
}

func accessTokenFromRequest(r *http.Request) string {
	if token := bearerToken(r.Header.Get("Authorization")); token != "" {
		return token
	}

	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func tokenSource(r *http.Request) string {
	if bearerToken(r.Header.Get("Authorization")) != "" {
		return "authorization_header"
	}

	if _, err := r.Cookie(accessTokenCookieName); err == nil {
		return "cookie"
	}

	return "none"
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: "Unauthorized",
	})
}
