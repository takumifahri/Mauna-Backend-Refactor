package service

import (
	"context"
	"log/slog"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/usecase"
)

type rateLimiterService struct {
	rateLimitRepo domain.RateLimitRepository
	now           func() time.Time
}

var _ usecase.RateLimiterUsecase = (*rateLimiterService)(nil)

func NewRateLimiterService(rateLimitRepo domain.RateLimitRepository, cleanupInterval time.Duration) usecase.RateLimiterUsecase {
	limiter := newRateLimiterService(rateLimitRepo, cleanupInterval, time.Now)
	return limiter
}

func newRateLimiterService(rateLimitRepo domain.RateLimitRepository, cleanupInterval time.Duration, now func() time.Time) *rateLimiterService {
	limiter := &rateLimiterService{
		rateLimitRepo: rateLimitRepo,
		now:           now,
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

	count, resetAt, err := s.rateLimitRepo.Increment(ctx, key, policy.Window)
	if err != nil {
		return usecase.RateLimitDecision{}, domain.NewInternalError(err)
	}
	if count > policy.Limit {
		return usecase.RateLimitDecision{
			Allowed:    false,
			RetryAfter: resetAt.Sub(s.now()),
		}, nil
	}

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
	if err := s.rateLimitRepo.DeleteExpired(context.Background(), s.now().Add(-time.Minute)); err != nil {
		slog.Warn("rate_limit_cleanup_failed", slog.Any("error", err))
	}
}
