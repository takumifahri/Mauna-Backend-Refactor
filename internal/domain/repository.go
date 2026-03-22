package domain

import (
    "context"

    "REFACTORING_MAUNA/internal/domain/entities"
)

// UserRepository interface untuk User operations
type UserRepository interface {
    // Create creates a new user
    Create(ctx context.Context, user *entities.User) (int64, error)

    // GetByID gets user by ID
    GetByID(ctx context.Context, id int64) (*entities.User, error)

    // GetByEmail gets user by email
    GetByEmail(ctx context.Context, email string) (*entities.User, error)

    // GetByUsername gets user by username
    GetByUsername(ctx context.Context, username string) (*entities.User, error)

    // GetByEmailOrUsername gets user by email or username
    GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (*entities.User, error)

    // GetAll gets all users (dengan pagination)
    GetAll(ctx context.Context, limit int, offset int) ([]entities.User, int64, error)

    // Update updates user data
    Update(ctx context.Context, user *entities.User) error

    // Delete deletes user (soft delete)
    Delete(ctx context.Context, id int64) error

    // CheckEmailExists checks if email exists
    CheckEmailExists(ctx context.Context, email string) (bool, error)

    // CheckUsernameExists checks if username exists
    CheckUsernameExists(ctx context.Context, username string) (bool, error)
}

// BadgeRepository interface untuk Badge operations
type BadgeRepository interface {
    Create(ctx context.Context, badge *entities.Badge) (int64, error)
    GetByID(ctx context.Context, id int64) (*entities.Badge, error)
    GetAll(ctx context.Context) ([]entities.Badge, error)
    Update(ctx context.Context, badge *entities.Badge) error
    Delete(ctx context.Context, id int64) error
}

// DictionaryRepository interface untuk Dictionary operations
type DictionaryRepository interface {
    Create(ctx context.Context, kamus *entities.Kamus) (int64, error)
    GetByID(ctx context.Context, id int64) (*entities.Kamus, error)
    GetAll(ctx context.Context, limit int, offset int) ([]entities.Kamus, int64, error)
    Search(ctx context.Context, keyword string, category string) ([]entities.Kamus, error)
    Update(ctx context.Context, kamus *entities.Kamus) error
    Delete(ctx context.Context, id int64) error
}

// LevelRepository interface untuk Level operations
type LevelRepository interface {
    Create(ctx context.Context, level *entities.Level) (int64, error)
    GetByID(ctx context.Context, id int64) (*entities.Level, error)
    GetAll(ctx context.Context) ([]entities.Level, error)
    Update(ctx context.Context, level *entities.Level) error
    Delete(ctx context.Context, id int64) error
}

// QuestionRepository interface untuk Question operations
type QuestionRepository interface {
    Create(ctx context.Context, soal *entities.Soal) (int64, error)
    GetByID(ctx context.Context, id int64) (*entities.Soal, error)
    GetAll(ctx context.Context, limit int, offset int) ([]entities.Soal, int64, error)
    GetBySublevels(ctx context.Context, sublevelID int64) ([]entities.Soal, error)
    Search(ctx context.Context, keyword string) ([]entities.Soal, error)
    Update(ctx context.Context, soal *entities.Soal) error
    Delete(ctx context.Context, id int64) error
}

// ProgressRepository interface untuk Progress operations
type ProgressRepository interface {
    Create(ctx context.Context, progress *entities.Progress) (int64, error)
    GetByUserID(ctx context.Context, userID int64) (*entities.Progress, error)
    Update(ctx context.Context, progress *entities.Progress) error
    Delete(ctx context.Context, id int64) error
}