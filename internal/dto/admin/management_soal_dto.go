package admin

import "time"

type CreateManagementSoalRequest struct {
	Question        string `json:"question"`
	Answer          string `json:"answer"`
	VideoURL        string `json:"video_url,omitempty"`
	ImageURL        string `json:"image_url,omitempty"`
	DictionaryID    int64  `json:"dictionary_id"`
	PointGamifikasi int    `json:"point_gamifikasi,omitempty"`
	SubLevelID      int64  `json:"sublevel_id"`
	Categories      string `json:"categories"`
}

type UpdateManagementSoalRequest struct {
	Question        *string `json:"question,omitempty"`
	Answer          *string `json:"answer,omitempty"`
	VideoURL        *string `json:"video_url,omitempty"`
	ImageURL        *string `json:"image_url,omitempty"`
	DictionaryID    *int64  `json:"dictionary_id,omitempty"`
	PointGamifikasi *int    `json:"point_gamifikasi,omitempty"`
	SubLevelID      *int64  `json:"sublevel_id,omitempty"`
	Categories      *string `json:"categories,omitempty"`
}

type ManagementSoalFilter struct {
	Query          string
	ID             int64
	SubLevelID     int64
	DictionaryID   int64
	Categories     string
	IncludeDeleted bool
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ManagementSoalListResponse struct {
	Soal       []ManagementSoalResponse `json:"soal"`
	Pagination PaginationResponse       `json:"pagination"`
}

type ManagementSoalResponse struct {
	ID              int64      `json:"id" db:"id"`
	Question        string     `json:"question" db:"question"`
	Answer          string     `json:"answer" db:"answer"`
	VideoURL        string     `json:"video_url,omitempty" db:"video_url"`
	ImageURL        string     `json:"image_url,omitempty" db:"image_url"`
	DictionaryID    int64      `json:"dictionary_id" db:"dictionary_id"`
	PointGamifikasi int        `json:"point_gamifikasi" db:"point_gamifikasi"`
	SubLevelID      int64      `json:"sublevel_id" db:"sublevel_id"`
	Categories      string     `json:"categories" db:"categories"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
