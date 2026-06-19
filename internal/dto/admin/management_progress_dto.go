package admin

import "time"

type ManagementProgressFilter struct {
	Query          string
	ID             int64
	UserID         string
	SubLevelID     int64
	Status         string
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementProgressListResponse struct {
	Progress   []ManagementProgressResponse `json:"progress"`
	Pagination PaginationResponse           `json:"pagination"`
}

type ManagementProgressResponse struct {
	ID                   int64      `json:"id" db:"id"`
	UserID               string     `json:"user_id" db:"user_id"`
	SubLevelID           int64      `json:"sublevel_id" db:"sublevel_id"`
	Status               string     `json:"status" db:"status"`
	TotalQuestions       int        `json:"total_questions" db:"total_questions"`
	CorrectAnswers       int        `json:"correct_answers" db:"correct_answers"`
	Score                int        `json:"score" db:"score"`
	Stars                int        `json:"stars" db:"stars"`
	CompletionPercentage int        `json:"completion_percentage" db:"completion_percentage"`
	Attempts             int        `json:"attempts" db:"attempts"`
	BestScore            int        `json:"best_score" db:"best_score"`
	BestStars            int        `json:"best_stars" db:"best_stars"`
	IsUnlocked           bool       `json:"is_unlocked" db:"is_unlocked"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	LastAttempt          *time.Time `json:"last_attempt,omitempty" db:"last_attempt"`
	CompletedAt          *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}
