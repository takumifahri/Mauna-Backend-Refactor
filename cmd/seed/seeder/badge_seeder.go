package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
    "REFACTORING_MAUNA/internal/domain/entities"
)

type BadgeSeeder struct {
    db *sqlx.DB
}

func NewBadgeSeeder(db *sqlx.DB) *BadgeSeeder {
    return &BadgeSeeder{db: db}
}

func (s *BadgeSeeder) Name() string {
    return "BadgeSeeder"
}

type badgeData struct {
    Nama      string
    Deskripsi string
    Icon      string
    Level     entities.DifficultyLevel
}

func (s *BadgeSeeder) Run() error {
    PrintInfo("Seeding badges...")

    badgesData := []badgeData{
        {
            Nama:      "First Steps",
            Deskripsi: "Complete your first lesson",
            Icon:      "badges/First_Steps.png",
            Level:     entities.DifficultyEasy,
        },
        {
            Nama:      "Alphabet Master",
            Deskripsi: "Master all alphabet signs",
            Icon:      "badges/Alphabet_Master.png",
            Level:     entities.DifficultyMedium,
        },
        {
            Nama:      "Number Expert",
            Deskripsi: "Perfect number signs 0-9",
            Icon:      "badges/Numbers_Expert.png",
            Level:     entities.DifficultyMedium,
        },
        {
            Nama:      "Conversation Pro",
            Deskripsi: "Complete advanced conversations",
            Icon:      "badges/Conversation.png",
            Level:     entities.DifficultyHard,
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
    for _, badgeData := range badgesData {
        // Check if badge already exists
        var existingID sql.NullInt64
        err := tx.QueryRow(
            "SELECT id FROM badges WHERE nama = $1",
            badgeData.Nama,
        ).Scan(&existingID)

        if err == nil && existingID.Valid {
            // Badge already exists
            PrintWarning(fmt.Sprintf("Badge already exists: %s", badgeData.Nama))
            continue
        }

        if err != nil && err != sql.ErrNoRows {
            // Unexpected error
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        // Insert new badge
        _, err = tx.Exec(
            "INSERT INTO badges (nama, deskripsi, icon, level) VALUES ($1, $2, $3, $4)",
            badgeData.Nama,
            badgeData.Deskripsi,
            badgeData.Icon,
            badgeData.Level,
        )

        if err != nil {
            PrintError(fmt.Sprintf("Failed to create badge: %v", err))
            return err
        }

        PrintSuccess(fmt.Sprintf("Created badge: %s", badgeData.Nama))
        createdCount++
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Badge seeding completed. Created %d new badges.", createdCount))
    return nil
}