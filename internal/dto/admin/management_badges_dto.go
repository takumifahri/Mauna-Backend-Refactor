package admin

import "time"

type CreateManagementBadgeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Level       string `json:"level"`
}

type UpdateManagementBadgeRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Level       *string `json:"level,omitempty"`
}

type ManagementBadgesFilter struct {
	Query          string
	ID             int64
	Name           string
	Level          string
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementBadgesListResponse struct {
	Badges     []ManagementBadgeResponse `json:"badges"`
	Pagination PaginationResponse        `json:"pagination"`
}

type ManagementBadgeResponse struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	Icon        string     `json:"icon,omitempty" db:"icon"`
	Level       string     `json:"level" db:"level"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
