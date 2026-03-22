package entities

import "time"

type DifficultyLevel string

const (
    DifficultyEasy   DifficultyLevel = "easy"
    DifficultyMedium DifficultyLevel = "medium"
    DifficultyHard   DifficultyLevel = "hard"
)

type Badge struct {
    ID        int64             `db:"id"`
    Nama      string            `db:"nama"`
    Deskripsi *string           `db:"deskripsi"`
    Icon      *string           `db:"icon"`
    Level     DifficultyLevel   `db:"level"`
    CreatedAt time.Time         `db:"created_at"`
    UpdatedAt *time.Time        `db:"updated_at"`
    DeletedAt *time.Time        `db:"deleted_at"`
}

type UserBadge struct {
    ID       int64     `db:"id"`
    UserID   int64     `db:"user_id"`
    BadgeID  int64     `db:"badge_id"`
    EarnedAt time.Time `db:"earned_at"`
}