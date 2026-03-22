package entities

import "time"

type TokenBlacklist struct {
    ID        int64      `db:"id"`
    Token     string     `db:"token"`
    UserID    int64      `db:"user_id"`
    RevokedAt time.Time  `db:"revoked_at"`
    ExpiresAt time.Time  `db:"expires_at"`
    Reason    *string    `db:"reason"`
}