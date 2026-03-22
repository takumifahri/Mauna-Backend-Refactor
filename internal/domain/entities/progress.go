package entities

import "time"

type ProgressStatus string

const (
    StatusNotStarted ProgressStatus = "not_started"
    StatusInProgress ProgressStatus = "in_progress"
    StatusCompleted  ProgressStatus = "completed"
    StatusFailed     ProgressStatus = "failed"
)

type Progress struct {
    ID                  int64              `db:"id"`
    UserID              int64              `db:"user_id"`
    SubLevelID          int64              `db:"sublevel_id"`
    Status              ProgressStatus     `db:"status"`
    TotalQuestions      int                `db:"total_questions"`
    CorrectAnswers      int                `db:"correct_answers"`
    Score               int                `db:"score"`
    Stars               int                `db:"stars"`
    CompletionPercentage int               `db:"completion_percentage"`
    Attempts            int                `db:"attempts"`
    BestScore           int                `db:"best_score"`
    BestStars           int                `db:"best_stars"`
    IsUnlocked          bool               `db:"is_unlocked"`
    CreatedAt           time.Time          `db:"created_at"`
    UpdatedAt           *time.Time         `db:"updated_at"`
    DeletedAt           *time.Time         `db:"deleted_at"`
    LastAttempt         *time.Time         `db:"last_attempt"`
    CompletedAt         *time.Time         `db:"completed_at"`
}