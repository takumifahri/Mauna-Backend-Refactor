package main

import (
    "fmt"
    "log"

    "REFACTORING_MAUNA/cmd/seed/seeder"
    "REFACTORING_MAUNA/config"
    "REFACTORING_MAUNA/pkg/database"
)

func main() {
    // Load database config
    dbConfig := config.NewDatabaseConfig()

    // Create database connection
    db, err := database.New(dbConfig)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    fmt.Println("🌱 Starting database seeding...")

    // Run seeders in order (respecting FK dependencies)
    seeders := []seeder.BaseSeeder{
        seeder.NewBadgeSeeder(db.DB),
        seeder.NewShopSeeder(db.DB),
        seeder.NewUserSeeder(db.DB),
        seeder.NewKamusSeeder(db.DB),
        seeder.NewLevelSeeder(db.DB),
        seeder.NewUserBadgeSeeder(db.DB),
        seeder.NewSublevelSeeder(db.DB),
        seeder.NewSoalSeeder(db.DB),
    }

    for _, s := range seeders {
        if err := s.Run(); err != nil {
            log.Fatalf("Seeder %s failed: %v", s.Name(), err)
        }
    }

    fmt.Println("✅ All seeders completed successfully!")
}