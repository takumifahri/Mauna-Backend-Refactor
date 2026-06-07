package service

import (
	"context"
	"testing"
	"time"

	"REFACTORING_MAUNA/internal/usecase"
)

type fakeRateLimitRepository struct {
	currentTime func() time.Time
	entries     map[string]fakeRateLimitEntry
}

type fakeRateLimitEntry struct {
	count   int
	resetAt time.Time
}

func newFakeRateLimitRepository(currentTime func() time.Time) *fakeRateLimitRepository {
	return &fakeRateLimitRepository{
		currentTime: currentTime,
		entries:     make(map[string]fakeRateLimitEntry),
	}
}

func (r *fakeRateLimitRepository) Increment(ctx context.Context, key string, window time.Duration) (int, time.Time, error) {
	now := r.currentTime()
	entry := r.entries[key]
	if entry.resetAt.IsZero() || !now.Before(entry.resetAt) {
		entry = fakeRateLimitEntry{
			count:   0,
			resetAt: now.Add(window),
		}
	}

	entry.count++
	r.entries[key] = entry
	return entry.count, entry.resetAt, nil
}

func (r *fakeRateLimitRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	for key, entry := range r.entries {
		if !entry.resetAt.After(before) {
			delete(r.entries, key)
		}
	}
	return nil
}

func TestRateLimiterServiceBlocksAfterLimit(t *testing.T) {
	currentTime := time.Date(2026, 5, 27, 10, 0, 0, 0, time.UTC)
	now := func() time.Time {
		return currentTime
	}
	limiter := newRateLimiterService(newFakeRateLimitRepository(now), 0, now)
	policy := usecase.RateLimitPolicy{Limit: 2, Window: time.Minute}

	for i := 0; i < 2; i++ {
		decision, err := limiter.Allow(context.Background(), "POST:/api/auth/login:192.0.2.10", policy)
		if err != nil {
			t.Fatalf("allow request %d error = %v", i+1, err)
		}
		if !decision.Allowed {
			t.Fatalf("request %d allowed = false, want true", i+1)
		}
	}

	decision, err := limiter.Allow(context.Background(), "POST:/api/auth/login:192.0.2.10", policy)
	if err != nil {
		t.Fatalf("blocked request error = %v", err)
	}
	if decision.Allowed {
		t.Fatal("blocked request allowed = true, want false")
	}
	if decision.RetryAfter <= 0 {
		t.Fatalf("retry after = %s, want positive duration", decision.RetryAfter)
	}
}

func TestRateLimiterServiceResetsAfterWindow(t *testing.T) {
	currentTime := time.Date(2026, 5, 27, 10, 0, 0, 0, time.UTC)
	now := func() time.Time {
		return currentTime
	}
	limiter := newRateLimiterService(newFakeRateLimitRepository(now), 0, now)
	policy := usecase.RateLimitPolicy{Limit: 1, Window: time.Minute}
	key := "POST:/api/auth/login:192.0.2.10"

	first, err := limiter.Allow(context.Background(), key, policy)
	if err != nil {
		t.Fatalf("first request error = %v", err)
	}
	if !first.Allowed {
		t.Fatal("first request allowed = false, want true")
	}

	blocked, err := limiter.Allow(context.Background(), key, policy)
	if err != nil {
		t.Fatalf("blocked request error = %v", err)
	}
	if blocked.Allowed {
		t.Fatal("blocked request allowed = true, want false")
	}

	currentTime = currentTime.Add(time.Minute + time.Second)
	allowed, err := limiter.Allow(context.Background(), key, policy)
	if err != nil {
		t.Fatalf("reset request error = %v", err)
	}
	if !allowed.Allowed {
		t.Fatal("reset request allowed = false, want true")
	}
}
