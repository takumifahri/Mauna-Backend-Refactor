package main

import (
    "log"
    "github.com/joho/godotenv"
    "REFACTORING_MAUNA/pkg/database"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }

    // Connect to database
    db, err := database.NewFromEnv()
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }
    defer db.Close()

    // Check health
    if err := db.Health(); err != nil {
        log.Fatalf("Database health check failed: %v", err)
    }

    log.Println("Database connected successfully!")
    log.Println(db.Stats())

    // Your application code here...
}