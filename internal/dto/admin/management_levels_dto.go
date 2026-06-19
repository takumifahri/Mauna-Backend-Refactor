package admin

import "time"

// Level DTOs

type CreateManagementLevelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Tujuan      string `json:"tujuan,omitempty"`
}

type UpdateManagementLevelRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Tujuan      *string `json:"tujuan,omitempty"`
}

type ManagementLevelsFilter struct {
	Query          string
	ID             int64
	Name           string
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementLevelsListResponse struct {
	Levels     []ManagementLevelResponse `json:"levels"`
	Pagination PaginationResponse        `json:"pagination"`
}

type ManagementLevelResponse struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	Tujuan      string     `json:"tujuan,omitempty" db:"tujuan"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// SubLevel DTOs

type CreateManagementSubLevelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Tujuan      string `json:"tujuan,omitempty"`
	LevelID     int64  `json:"level_id"`
}

type UpdateManagementSubLevelRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Tujuan      *string `json:"tujuan,omitempty"`
	LevelID     *int64  `json:"level_id,omitempty"`
}

type ManagementSubLevelsFilter struct {
	Query          string
	ID             int64
	Name           string
	LevelID        int64
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementSubLevelsListResponse struct {
	SubLevels  []ManagementSubLevelResponse `json:"sub_levels"`
	Pagination PaginationResponse           `json:"pagination"`
}

type ManagementSubLevelResponse struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	Tujuan      string     `json:"tujuan,omitempty" db:"tujuan"`
	LevelID     int64      `json:"level_id" db:"level_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
