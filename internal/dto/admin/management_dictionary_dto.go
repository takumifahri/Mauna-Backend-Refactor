package admin

import "time"

type CreateManagementDictionaryRequest struct {
	WordText    string `json:"word_text"`
	Definition  string `json:"definition"`
	Category    string `json:"category"`
	ImageURLRef string `json:"image_url_ref,omitempty"`
	VideoURL    string `json:"video_url,omitempty"`
}

type UpdateManagementDictionaryRequest struct {
	WordText    *string `json:"word_text,omitempty"`
	Definition  *string `json:"definition,omitempty"`
	Category    *string `json:"category,omitempty"`
	ImageURLRef *string `json:"image_url_ref,omitempty"`
	VideoURL    *string `json:"video_url,omitempty"`
}

type ManagementDictionaryFilter struct {
	Query          string
	ID             int64
	WordText       string
	Category       string
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementDictionaryListResponse struct {
	Dictionaries []ManagementDictionaryResponse `json:"dictionaries"`
	Pagination   PaginationResponse             `json:"pagination"`
}

type ManagementDictionaryResponse struct {
	ID          int64      `json:"id" db:"id"`
	WordText    string     `json:"word_text" db:"word_text"`
	Definition  string     `json:"definition" db:"definition"`
	Category    string     `json:"category" db:"category"`
	ImageURLRef string     `json:"image_url_ref,omitempty" db:"image_url_ref"`
	VideoURL    string     `json:"video_url,omitempty" db:"video_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
