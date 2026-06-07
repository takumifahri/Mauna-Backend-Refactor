package repository

import (
	"context"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/pkg/database"
)

type rateLimitRepository struct {
	db *database.DB
}

func NewRateLimitRepository(db *database.DB) domain.RateLimitRepository {
	return &rateLimitRepository{db: db}
}

func (r *rateLimitRepository) Increment(ctx context.Context, key string, window time.Duration) (int, time.Time, error) {
	query := `INSERT INTO rate_limits (key, count, reset_at, updated_at)
             VALUES ($1, 1, NOW() + make_interval(secs => $2), NOW())
             ON CONFLICT (key) DO UPDATE SET
                 count = CASE
                     WHEN rate_limits.reset_at <= NOW() THEN 1
                     ELSE rate_limits.count + 1
                 END,
                 reset_at = CASE
                     WHEN rate_limits.reset_at <= NOW() THEN NOW() + make_interval(secs => $2)
                     ELSE rate_limits.reset_at
                 END,
                 updated_at = NOW()
             RETURNING count, reset_at`

	var count int
	var resetAt time.Time
	err := r.db.QueryRowContext(ctx, query, key, int(window.Seconds())).Scan(&count, &resetAt)
	return count, resetAt, err
}

func (r *rateLimitRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := `DELETE FROM rate_limits WHERE reset_at <= $1`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}
