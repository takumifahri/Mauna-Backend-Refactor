package entities

import "time"

type TokenBlacklist struct {
	ID        int64      `db:"id"`
	Token     string     `db:"token"`
	UserID    string     `db:"user_id"`
	RevokedAt time.Time  `db:"revoked_at"`
	ExpiresAt time.Time  `db:"expires_at"`
	Reason    *string    `db:"reason"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
