package entities

import "time"

type DailyTask struct {
    ID                  int64      `db:"id"`
    UserID              int64      `db:"user_id"`
    Date                *time.Time `db:"date"`
    CompletedSublevels  int        `db:"completed_sublevels"`
    IsCompleted         bool       `db:"is_completed"`
    LastUpdate          time.Time  `db:"last_update"`
    CreatedAt           time.Time  `db:"created_at"`
    UpdatedAt           *time.Time `db:"updated_at"`
    DeletedAt           *time.Time `db:"deleted_at"`
}