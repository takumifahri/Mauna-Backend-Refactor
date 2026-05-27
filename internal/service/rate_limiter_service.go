package service

import (
	"context"
	"sync"
	"time"

	"REFACTORING_MAUNA/internal/usecase"
)

type rateLimitEntry struct {
	count     int
	resetAt   time.Time
	updatedAt time.Time
}

type rateLimiterService struct {
	mu      sync.Mutex
	entries map[string]rateLimitEntry
	now     func() time.Time
}

var _ usecase.RateLimiterUsecase = (*rateLimiterService)(nil)

func NewRateLimiterService(cleanupInterval time.Duration) usecase.RateLimiterUsecase {
	limiter := newRateLimiterService(cleanupInterval, time.Now)
	return limiter
}

func newRateLimiterService(cleanupInterval time.Duration, now func() time.Time) *rateLimiterService {
	limiter := &rateLimiterService{
		entries: make(map[string]rateLimitEntry),
		now:     now,
	}

	if cleanupInterval > 0 {
		go limiter.cleanupLoop(cleanupInterval)
	}

	return limiter
}

func (s *rateLimiterService) Allow(ctx context.Context, key string, policy usecase.RateLimitPolicy) (usecase.RateLimitDecision, error) {
	if policy.Limit <= 0 || policy.Window <= 0 {
		return usecase.RateLimitDecision{Allowed: true}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	entry := s.entries[key]
	if entry.resetAt.IsZero() || !now.Before(entry.resetAt) {
		entry = rateLimitEntry{
			count:     0,
			resetAt:   now.Add(policy.Window),
			updatedAt: now,
		}
	}

	entry.updatedAt = now
	if entry.count >= policy.Limit {
		s.entries[key] = entry
		return usecase.RateLimitDecision{
			Allowed:    false,
			RetryAfter: entry.resetAt.Sub(now),
		}, nil
	}

	entry.count++
	s.entries[key] = entry
	return usecase.RateLimitDecision{Allowed: true}, nil
}

func (s *rateLimiterService) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanup()
	}
}

func (s *rateLimiterService) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	for key, entry := range s.entries {
		if now.After(entry.resetAt) && now.Sub(entry.updatedAt) > time.Minute {
			delete(s.entries, key)
		}
	}
}
