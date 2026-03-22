package seeder

import (
    "fmt"
    "github.com/jmoiron/sqlx"
)

// BaseSeeder adalah interface untuk semua seeder
type BaseSeeder interface {
    Run() error
    Name() string
}

// BaseSeed adalah struct helper untuk database operations
type BaseSeed struct {
    DB *sqlx.DB
}

// NewBaseSeed membuat instance baru BaseSeed
func NewBaseSeed(db *sqlx.DB) *BaseSeed {
    return &BaseSeed{DB: db}
}

// PrintSuccess print pesan sukses dengan emoji
func PrintSuccess(message string) {
    fmt.Printf("✅ %s\n", message)
}

// PrintError print pesan error dengan emoji
func PrintError(message string) {
    fmt.Printf("❌ %s\n", message)
}

// PrintWarning print pesan warning dengan emoji
func PrintWarning(message string) {
    fmt.Printf("⚠️ %s\n", message)
}

// PrintInfo print pesan info dengan emoji
func PrintInfo(message string) {
    fmt.Printf("🌱 %s\n", message)
}

// ptrStr converts string to pointer
func ptrStr(s string) *string {
    return &s
}