package usecase

import (
	"context"
	"time"
)

type RateLimitPolicy struct {
	Limit  int
	Window time.Duration
}

type RateLimitDecision struct {
	Allowed    bool
	RetryAfter time.Duration
}

type RateLimiterUsecase interface {
	Allow(ctx context.Context, key string, policy RateLimitPolicy) (RateLimitDecision, error)
}
