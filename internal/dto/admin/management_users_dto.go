package admin

import (
	"time"
)

type CreateManagementUserRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Phone      string `json:"phone,omitempty"`
	Role       string `json:"role"`
	IsActive   *bool  `json:"is_active,omitempty"`
	IsVerified *bool  `json:"is_verified,omitempty"`
}

type UpdateManagementUserRequest struct {
	Username   *string `json:"username,omitempty"`
	Email      *string `json:"email,omitempty"`
	Password   *string `json:"password,omitempty"`
	Name       *string `json:"name,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Role       *string `json:"role,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	IsVerified *bool   `json:"is_verified,omitempty"`
}

type ManagementUsersFilter struct {
	Query          string
	Role           string
	IsActive       *bool
	IsVerified     *bool
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementUsersListResponse struct {
	Users      []ManagementUserResponse `json:"users"`
	Pagination PaginationResponse       `json:"pagination"`
}

type PaginationResponse struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type ManagementUserResponse struct {
	ID                    string     `json:"id"`
	Username              string     `json:"username"`
	Email                 string     `json:"email"`
	Name                  string     `json:"name"`
	Phone                 string     `json:"phone,omitempty"`
	Role                  string     `json:"role"`
	IsActive              bool       `json:"is_active"`
	IsVerified            bool       `json:"is_verified"`
	TotalBadges           int        `json:"total_badges"`
	TotalXP               int        `json:"total_xp"`
	TotalQuizzesCompleted int        `json:"total_quizzes_completed"`
	TotalPoints           int        `json:"total_points"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
	LastLogin             *time.Time `json:"last_login,omitempty"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty"`
}
