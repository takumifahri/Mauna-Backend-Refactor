package repository

import (
    "context"
    "database/sql"

    "REFACTORING_MAUNA/internal/domain"
    "REFACTORING_MAUNA/internal/domain/entities"
    "REFACTORING_MAUNA/pkg/database"
)

type userRepository struct {
    db *database.DB
}

func NewUserRepository(db *database.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) (int64, error) {
    query := `INSERT INTO users (username, email, password_hash, name, role, is_active, is_verified, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
             RETURNING id`

    var id int64
    err := r.db.QueryRowContext(ctx, query, 
        user.Username, 
        user.Email, 
        user.PasswordHash,  // ← ADD THIS (missing!)
        user.Nama,
        string(user.Role),  // Convert UserRole to string
        user.IsActive, 
        user.IsVerified,
    ).Scan(&id)
    return id, err
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entities.User, error) {
    query := `SELECT id, username, email, password_hash, name, role, is_active, is_verified, created_at, updated_at
             FROM users WHERE id = $1 AND deleted_at IS NULL`

    user := &entities.User{}
    err := r.db.GetContext(ctx, user, query, id)
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
    query := `SELECT id, username, email, password_hash, name, role, is_active, is_verified, created_at, updated_at
             FROM users WHERE email = $1 AND deleted_at IS NULL`

    user := &entities.User{}
    err := r.db.GetContext(ctx, user, query, email)
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    return user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
    query := `SELECT id, username, email, password_hash, name, role, is_active, is_verified, created_at, updated_at
             FROM users WHERE username = $1 AND deleted_at IS NULL`

    user := &entities.User{}
    err := r.db.GetContext(ctx, user, query, username)
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    return user, err
}

func (r *userRepository) GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (*entities.User, error) {
    query := `SELECT id, username, email, password_hash, name, role, is_active, is_verified, created_at, updated_at
             FROM users WHERE (email = $1 OR username = $1) AND deleted_at IS NULL`

    user := &entities.User{}
    err := r.db.GetContext(ctx, user, query, emailOrUsername)
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    return user, err
}

func (r *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

    var exists bool
    err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
    return exists, err
}

func (r *userRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL)`

    var exists bool
    err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
    return exists, err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
    query := `UPDATE users SET username = $1, email = $2, password_hash = $3, name = $4, 
             role = $5, is_active = $6, is_verified = $7, updated_at = NOW()
             WHERE id = $8 AND deleted_at IS NULL`

    result, err := r.db.ExecContext(ctx, query, 
        user.Username, 
        user.Email, 
        user.PasswordHash,  // ← ADD THIS (missing!)
        user.Nama,
        string(user.Role),  // Convert UserRole to string
        user.IsActive, 
        user.IsVerified, 
        user.ID,
    )
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }

    return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
    query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`

    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }

    return nil
}

func (r *userRepository) GetAll(ctx context.Context, limit int, offset int) ([]entities.User, int64, error) {
    query := `SELECT id, username, email, password_hash, name, role, is_active, is_verified, created_at, updated_at
             FROM users WHERE deleted_at IS NULL
             ORDER BY created_at DESC LIMIT $1 OFFSET $2`

    users := []entities.User{}
    err := r.db.SelectContext(ctx, &users, query, limit, offset)
    if err != nil {
        return nil, 0, err
    }

    // Get total count
    countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
    var total int64
    err = r.db.GetContext(ctx, &total, countQuery)
    if err != nil {
        return nil, 0, err
    }

    return users, total, nil
}