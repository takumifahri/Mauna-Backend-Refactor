package dto

import "time"

type ProfileResponse struct {
	ID                    int64      `json:"id"`
	UniqueID              string     `json:"unique_id"`
	Username              string     `json:"username"`
	Email                 string     `json:"email"`
	Name                  string     `json:"name"`
	Phone                 string     `json:"phone,omitempty"`
	Role                  string     `json:"role"`
	IsActive              bool       `json:"is_active"`
	IsVerified            bool       `json:"is_verified"`
	Avatar                string     `json:"avatar,omitempty"`
	Bio                   string     `json:"bio,omitempty"`
	TotalBadges           int        `json:"total_badges"`
	CurrentStreak         int        `json:"current_streak"`
	LongestStreak         int        `json:"longest_streak"`
	WeeklyXP              int        `json:"weekly_xp"`
	Tier                  string     `json:"tier"`
	TotalXP               int        `json:"total_xp"`
	TotalQuizzesCompleted int        `json:"total_quizzes_completed"`
	TotalPoints           int        `json:"total_points"`
	CreatedAt             time.Time  `json:"created_at"`
	LastLogin             *time.Time `json:"last_login,omitempty"`
}

type UpdateProfileRequest struct {
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	Name      *string `json:"name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
}
