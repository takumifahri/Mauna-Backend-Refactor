package entities

import "time"

type SoalType string

const (
    SoalTebakanGambar   SoalType = "TEBAK_GAMBAR"
    SoalOpenCamera      SoalType = "OPEN_CAMERA"
    SoalPilihanGanda    SoalType = "PILIHAN_GANDA"
    SoalMatematika      SoalType = "MATEMATIKA"
)

type Soal struct {
    ID                  int64      `db:"id"`
    Question            string     `db:"question"`
    Answer              string     `db:"answer"`
    VideoURL            *string    `db:"video_url"`
    ImageURL            *string    `db:"image_url"`
    DictionaryID        int64      `db:"dictionary_id"`
    PointGamifikasi     int        `db:"point_gamifikasi"`
    SubLevelID          int64      `db:"sublevel_id"`
    Categories          SoalType   `db:"categories"`
    CreatedAt           time.Time  `db:"created_at"`
    UpdatedAt           *time.Time `db:"updated_at"`
    DeletedAt           *time.Time `db:"deleted_at"`
}