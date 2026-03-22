package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
    "REFACTORING_MAUNA/internal/domain/entities"
    "REFACTORING_MAUNA/pkg/security"
    "time"
)

type UserSeeder struct {
    db *sqlx.DB
}

func NewUserSeeder(db *sqlx.DB) *UserSeeder {
    return &UserSeeder{db: db}
}

func (s *UserSeeder) Name() string {
    return "UserSeeder"
}

type userData struct {
    Username   string
    Email      string
    Password   string
    Nama       string
    Telpon     string
    Role       entities.UserRole
    IsActive   bool
    IsVerified bool
    Bio        string
}

func (s *UserSeeder) Run() error {
    PrintInfo("Seeding users...")

    usersData := []userData{
        {
            Username:   "admin",
            Email:      "admin@example.com",
            Password:   "AdminPass123!",
            Nama:       "Site Admin",
            Telpon:     "123-456-7890",
            Role:       entities.RoleAdmin,
            IsActive:   true,
            IsVerified: true,
            Bio:        "Administrator account",
        },
        {
            Username:   "moderator",
            Email:      "moderator@example.com",
            Password:   "ModPass123!",
            Nama:       "Site Moderator",
            Telpon:     "098-765-4321",
            Role:       entities.RoleModerator,
            IsActive:   true,
            IsVerified: true,
            Bio:        "Moderator account",
        },
        {
            Username:   "johndoe",
            Email:      "john.doe@example.com",
            Password:   "Password123",
            Nama:       "John Doe",
            Telpon:     "123-456-7890",
            Role:       entities.RoleUser,
            IsActive:   true,
            IsVerified: false,
            Bio:        "Sample user account",
        },
        {
            Username:   "janedoe",
            Email:      "jane.doe@example.com",
            Password:   "Password123",
            Nama:       "Jane Doe",
            Telpon:     "321-654-0987",
            Role:       entities.RoleUser,
            IsActive:   true,
            IsVerified: true,
            Bio:        "Another sample user",
        },
    }

    // Start transaction
    tx, err := s.db.Beginx()
    if err != nil {
        PrintError(fmt.Sprintf("Failed to start transaction: %v", err))
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    createdCount := 0
    for i, user := range usersData {
        // Check if user already exists by email or username
        var existingID sql.NullInt64
        err := tx.QueryRow(
            "SELECT id FROM users WHERE email = $1 OR username = $2",
            user.Email,
            user.Username,
        ).Scan(&existingID)

        if err == nil && existingID.Valid {
            PrintWarning(fmt.Sprintf("User already exists: %s (%s)", user.Username, user.Email))
            continue
        }

        if err != nil && err != sql.ErrNoRows {
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        // Hash password
        hashedPassword, err := security.HashPassword(user.Password)
        if err != nil {
            PrintError(fmt.Sprintf("Failed to hash password: %v", err))
            return err
        }

        // Generate unique_id (format: USR-00001)
        uniqueID := fmt.Sprintf("USR-%05d", i+1)

        // Insert user
        _, err = tx.Exec(
            `INSERT INTO users (
                unique_id, username, email, password, nama, telpon, role, 
                is_active, is_verified, bio, created_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
            uniqueID,
            user.Username,
            user.Email,
            hashedPassword,
            user.Nama,
            user.Telpon,
            user.Role,
            user.IsActive,
            user.IsVerified,
            user.Bio,
            time.Now(),
        )

        if err != nil {
            PrintError(fmt.Sprintf("Failed to create user: %v", err))
            return err
        }

        PrintSuccess(fmt.Sprintf("Created user: %s (%s) - ID: %s", user.Username, user.Email, uniqueID))
        createdCount++
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("User seeding completed. Created %d new users.", createdCount))
    return nil
}