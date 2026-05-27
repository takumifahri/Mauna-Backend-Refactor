package entities

import (
	"time"

	"REFACTORING_MAUNA/internal/domain/constant"
)



type Kamus struct {
	ID          int64              `db:"id"`
	WordText    string             `db:"word_text"`
	Definition  string             `db:"definition"`
	Category    constant.DictionaryCategory `db:"category"`
	ImageURLRef *string            `db:"image_url_ref"`
	VideoURL    *string            `db:"video_url"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   *time.Time         `db:"updated_at"`
	DeletedAt   *time.Time         `db:"deleted_at"`
}
