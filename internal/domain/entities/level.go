package entities

import "time"

type Level struct {
    ID          int64      `db:"id"`
    Name        string     `db:"name"`
    Description *string    `db:"description"`
    Tujuan      *string    `db:"tujuan"`
    CreatedAt   time.Time  `db:"created_at"`
    UpdatedAt   *time.Time `db:"updated_at"`
    DeletedAt   *time.Time `db:"deleted_at"`
}

type SubLevel struct {
    ID          int64      `db:"id"`
    Name        string     `db:"name"`
    Description *string    `db:"description"`
    Tujuan      *string    `db:"tujuan"`
    LevelID     int64      `db:"level_id"`
    CreatedAt   time.Time  `db:"created_at"`
    UpdatedAt   *time.Time `db:"updated_at"`
    DeletedAt   *time.Time `db:"deleted_at"`
}