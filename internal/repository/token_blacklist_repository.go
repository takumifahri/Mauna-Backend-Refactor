package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/pkg/database"
)

type tokenBlacklistRepository struct {
	db *database.DB
}

func NewTokenBlacklistRepository(db *database.DB) domain.TokenBlacklistRepository {
	return &tokenBlacklistRepository{db: db}
}

func (r *tokenBlacklistRepository) Create(ctx context.Context, token *entities.TokenBlacklist) error {
	query := `INSERT INTO token_blacklist (token, user_id, expires_at, reason)
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (token) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, token.Token, token.UserID, token.ExpiresAt, token.Reason)
	return err
}

func (r *tokenBlacklistRepository) Exists(ctx context.Context, token string) (bool, error) {
	query := `SELECT EXISTS(
                 SELECT 1 FROM token_blacklist
                 WHERE token = $1 AND expires_at > NOW()
             )`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, token).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return exists, err
}

func (r *tokenBlacklistRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := `DELETE FROM token_blacklist WHERE expires_at <= $1`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}
