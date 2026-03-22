package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
)

type LevelSeeder struct {
    db *sqlx.DB
}

func NewLevelSeeder(db *sqlx.DB) *LevelSeeder {
    return &LevelSeeder{db: db}
}

func (s *LevelSeeder) Name() string {
    return "LevelSeeder"
}

type levelData struct {
    Name        string
    Description string
    Tujuan      string
}

func (s *LevelSeeder) Run() error {
    PrintInfo("📚 Seeding Levels...")

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

    levelsData := []levelData{
        {Name: "Level 1", Description: "Pengenalan Abjad", Tujuan: "Mengenal dan memahami bahasa isyarat alphabet A-Z"},
        {Name: "Level 2", Description: "Pengenalan Kata", Tujuan: "Mengenal dan memahami bahasa isyarat kosakata dasar sehari-hari"},
        {Name: "Level 3", Description: "Angka dan Matematika dasar", Tujuan: "Mengenal dan memahami bahasa isyarat matematika dasar (penjumlahan dan pengurangan)"},
        {Name: "Level 4", Description: "Pengenalan Kata", Tujuan: "Mengenal dan memahami bahasa isyarat kosakata lanjutan"},
    }

    createdCount := 0
    for _, data := range levelsData {
        var existingID sql.NullInt64
        err := tx.QueryRow("SELECT id FROM level WHERE name = $1", data.Name).Scan(&existingID)

        if err == nil && existingID.Valid {
            PrintWarning(fmt.Sprintf("Level already exists: %s", data.Name))
            continue
        }

        if err != nil && err != sql.ErrNoRows {
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        _, err = tx.Exec(
            "INSERT INTO level (name, description, tujuan, created_at) VALUES ($1, $2, $3, NOW())",
            data.Name,
            data.Description,
            data.Tujuan,
        )

        if err != nil {
            PrintError(fmt.Sprintf("Failed to create level: %v", err))
            return err
        }

        PrintSuccess(fmt.Sprintf("Created level: %s", data.Name))
        createdCount++
    }

    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Level seeding completed. Created %d levels.", createdCount))
    return nil
}