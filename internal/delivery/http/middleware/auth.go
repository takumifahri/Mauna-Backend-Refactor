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
	return Auth(tokenManager)(next)
}

func Auth(tokenManager usecase.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, source := accessTokenFromRequest(r)
			if token == "" {
				logRequestError(r, "auth_token_missing", http.StatusUnauthorized, nil)
				writeUnauthorized(w)
				return
			}

			claims, err := tokenManager.VerifyToken(token)
			if err != nil {
				logRequestError(r, "auth_token_invalid", http.StatusUnauthorized, err, slog.String("token_source", source))
				writeUnauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), authClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	allowedRoles := map[string]struct{}{}
	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role != "" {
			allowedRoles[role] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				logRequestError(r, "auth_claims_missing", http.StatusUnauthorized, nil)
				writeUnauthorized(w)
				return
			}

			role := strings.ToLower(strings.TrimSpace(claims.Role))
			if _, ok := allowedRoles[role]; !ok {
				logRequestError(r, "auth_role_forbidden", http.StatusForbidden, nil, slog.String("role", claims.Role))
				writeForbidden(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ClaimsFromContext(ctx context.Context) (usecase.TokenClaims, bool) {
	claims, ok := ctx.Value(authClaimsKey).(usecase.TokenClaims)
	return claims, ok
}

func accessTokenFromRequest(r *http.Request) (string, string) {
	if token := bearerToken(r.Header.Get("Authorization")); token != "" {
		return token, "authorization_header"
	}

	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return "", "none"
	}
	return cookie.Value, "cookie"
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

func writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: "Forbidden",
		Error:   "forbidden",
	})
}
