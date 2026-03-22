package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
    "REFACTORING_MAUNA/internal/domain/entities"
)

type ShopSeeder struct {
    db *sqlx.DB
}

func NewShopSeeder(db *sqlx.DB) *ShopSeeder {
    return &ShopSeeder{db: db}
}

func (s *ShopSeeder) Name() string {
    return "ShopSeeder"
}

type shopItemData struct {
    Name        string
    Description string
    ItemType    entities.ShopItemType
    XpCost      int
    Icon        string
}

func (s *ShopSeeder) Run() error {
    PrintInfo("Starting shop item seeding...")

    itemsData := []shopItemData{
        {
            Name:        "Streak Freeze",
            Description: "Melindungi streak harian agar tidak hilang.",
            ItemType:    entities.TypeStreakFreeze,
            XpCost:      100,
            Icon:        "streak_freeze.png",
        },
        {
            Name:        "Badge Emas",
            Description: "Badge spesial untuk pencapaian emas.",
            ItemType:    entities.TypeBadge,
            XpCost:      250,
            Icon:        "badge_emas.png",
        },
        {
            Name:        "XP Boost",
            Description: "Dapatkan XP tambahan selama 1 jam.",
            ItemType:    entities.TypeBoost,
            XpCost:      150,
            Icon:        "xp_boost.png",
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

    PrintInfo("Seeding Shop Items...")
    createdCount := 0

    for _, item := range itemsData {
        // Check if shop item already exists
        var existingID sql.NullInt64
        err := tx.QueryRow(
            "SELECT id FROM shop_items WHERE name = $1",
            item.Name,
        ).Scan(&existingID)

        if err == nil && existingID.Valid {
            // Item already exists
            PrintWarning(fmt.Sprintf("Shop item already exists: %s", item.Name))
            continue
        }

        if err != nil && err != sql.ErrNoRows {
            // Unexpected error
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        // Insert new shop item
        _, err = tx.Exec(
            "INSERT INTO shop_items (name, description, item_type, xp_cost, icon) VALUES ($1, $2, $3, $4, $5)",
            item.Name,
            item.Description,
            item.ItemType,
            item.XpCost,
            item.Icon,
        )

        if err != nil {
            PrintError(fmt.Sprintf("Failed to create shop item: %v", err))
            return err
        }

        PrintSuccess(fmt.Sprintf("Created shop item: %s", item.Name))
        createdCount++
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Shop item seeding completed. Created %d items.", createdCount))
    return nil
}