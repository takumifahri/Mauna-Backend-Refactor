package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"REFACTORING_MAUNA/internal/usecase"
)

type fakeTokenManager struct {
	claims usecase.TokenClaims
	err    error
	token  string
}

func (m *fakeTokenManager) GenerateAccessToken(userID string, username, email, role string) (string, error) {
	return "", nil
}

func (m *fakeTokenManager) GenerateRefreshToken(userID string) (string, error) {
	return "", nil
}

func (m *fakeTokenManager) VerifyToken(token string) (usecase.TokenClaims, error) {
	m.token = token
	if m.err != nil {
		return usecase.TokenClaims{}, m.err
	}

	return m.claims, nil
}

func (m *fakeTokenManager) GetTokenExpiry(token string) (time.Time, error) {
	return time.Time{}, nil
}

type fakeRateLimiter struct {
	key string
}

func (l *fakeRateLimiter) Allow(ctx context.Context, key string, policy usecase.RateLimitPolicy) (usecase.RateLimitDecision, error) {
	l.key = key
	return usecase.RateLimitDecision{Allowed: true}, nil
}

func TestChainRunsMiddlewareInDeclaredOrder(t *testing.T) {
	var calls []string
	first := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "first_before")
			next.ServeHTTP(w, r)
			calls = append(calls, "first_after")
		})
	}
	second := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "second_before")
			next.ServeHTTP(w, r)
			calls = append(calls, "second_after")
		})
	}

	handler := Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "handler")
	}), first, second)

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	want := []string{"first_before", "second_before", "handler", "second_after", "first_after"}
	if len(calls) != len(want) {
		t.Fatalf("calls length = %d, want %d (%v)", len(calls), len(want), calls)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("calls[%d] = %q, want %q (%v)", i, calls[i], want[i], calls)
		}
	}
}

func TestAuthStoresClaimsAndPrefersBearerToken(t *testing.T) {
	tokenManager := &fakeTokenManager{
		claims: usecase.TokenClaims{UserID: "user-1", Role: "admin"},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	req.AddCookie(&http.Cookie{Name: accessTokenCookieName, Value: "cookie-token"})

	handler := Auth(tokenManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("claims were not stored in context")
		}
		if claims.UserID != "user-1" || claims.Role != "admin" {
			t.Fatalf("claims = %+v, want user-1 admin", claims)
		}
	}))

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if tokenManager.token != "header-token" {
		t.Fatalf("verified token = %q, want header-token", tokenManager.token)
	}
}

func TestAuthRejectsInvalidToken(t *testing.T) {
	tokenManager := &fakeTokenManager{err: errors.New("bad token")}
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	rec := httptest.NewRecorder()

	Auth(tokenManager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRateLimitKeyUsesRoutePattern(t *testing.T) {
	limiter := &fakeRateLimiter{}
	middleware := NewRateLimitMiddleware(limiter)
	mux := http.NewServeMux()
	mux.Handle("GET /api/admin/users/{id}", middleware.Limiter(usecase.RateLimitPolicy{Limit: 1, Window: time.Minute})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/123", nil)
	req.RemoteAddr = "192.0.2.10:1234"
	mux.ServeHTTP(httptest.NewRecorder(), req)

	want := "GET:/api/admin/users/{id}:192.0.2.10"
	if limiter.key != want {
		t.Fatalf("key = %q, want %q", limiter.key, want)
	}
}
