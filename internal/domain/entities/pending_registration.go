package entities

import "time"

type PendingRegistration struct {
	ID           int64      `db:"id"`
	Username     string     `db:"username"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Nama         string     `db:"nama"`
	OTPHash      string     `db:"otp_hash"`
	Attempts     int        `db:"attempts"`
	ExpiresAt    time.Time  `db:"expires_at"`
	ConsumedAt   *time.Time `db:"consumed_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
