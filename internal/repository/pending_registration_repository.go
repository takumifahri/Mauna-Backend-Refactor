package repository

import (
	"context"
	"database/sql"
	"time"

	"REFACTORING_MAUNA/internal/domain"
	"REFACTORING_MAUNA/internal/domain/entities"
	"REFACTORING_MAUNA/pkg/database"
)

type pendingRegistrationRepository struct {
	db *database.DB
}

func NewPendingRegistrationRepository(db *database.DB) domain.PendingRegistrationRepository {
	return &pendingRegistrationRepository{db: db}
}

func (r *pendingRegistrationRepository) Upsert(ctx context.Context, registration *entities.PendingRegistration) error {
	query := `INSERT INTO pending_registrations (username, email, password_hash, nama, otp_hash, expires_at)
             VALUES ($1, $2, $3, $4, $5, $6)
             ON CONFLICT (email) DO UPDATE SET
                 username = EXCLUDED.username,
                 password_hash = EXCLUDED.password_hash,
                 nama = EXCLUDED.nama,
                 otp_hash = EXCLUDED.otp_hash,
                 attempts = 0,
                 expires_at = EXCLUDED.expires_at,
                 consumed_at = NULL,
                 updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query,
		registration.Username,
		registration.Email,
		registration.PasswordHash,
		registration.Nama,
		registration.OTPHash,
		registration.ExpiresAt,
	)
	return err
}

func (r *pendingRegistrationRepository) GetValidByEmail(ctx context.Context, email string) (*entities.PendingRegistration, error) {
	query := `SELECT id, username, email, password_hash, nama, otp_hash, attempts, expires_at, consumed_at, created_at, updated_at
             FROM pending_registrations
             WHERE email = $1 AND consumed_at IS NULL AND expires_at > NOW()`

	registration := &entities.PendingRegistration{}
	err := r.db.GetContext(ctx, registration, query, email)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidRegistrationOTP
	}
	return registration, err
}

func (r *pendingRegistrationRepository) IncrementAttempts(ctx context.Context, id int64) error {
	query := `UPDATE pending_registrations SET attempts = attempts + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *pendingRegistrationRepository) MarkConsumed(ctx context.Context, id int64) error {
	query := `UPDATE pending_registrations SET consumed_at = NOW(), updated_at = NOW()
             WHERE id = $1 AND consumed_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrInvalidRegistrationOTP
	}
	return nil
}

func (r *pendingRegistrationRepository) DeleteByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM pending_registrations WHERE email = $1`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

func (r *pendingRegistrationRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := `DELETE FROM pending_registrations WHERE expires_at <= $1 OR consumed_at IS NOT NULL`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}
