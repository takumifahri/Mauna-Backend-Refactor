package repository

import (
	"context"
	"database/sql"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/pkg/database"
)

type passwordResetTokenRepository struct {
	db *database.DB
}

func NewPasswordResetTokenRepository(db *database.DB) domain.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *entities.PasswordResetToken) error {
	query := `INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
             VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (r *passwordResetTokenRepository) GetValidByHash(ctx context.Context, tokenHash string) (*entities.PasswordResetToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, used_at, created_at
             FROM password_reset_tokens
             WHERE token_hash = $1 AND used_at IS NULL AND expires_at > NOW()`

	token := &entities.PasswordResetToken{}
	err := r.db.GetContext(ctx, token, query, tokenHash)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidResetToken
	}
	return token, err
}

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, id int64) error {
	query := `UPDATE password_reset_tokens SET used_at = NOW() WHERE id = $1 AND used_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrInvalidResetToken
	}

	return nil
}

func (r *passwordResetTokenRepository) DeleteActiveByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1 AND used_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *passwordResetTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at <= $1 OR used_at IS NOT NULL`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}
