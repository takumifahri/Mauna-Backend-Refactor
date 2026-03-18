package database

import (
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

// DB wraps database instance
type DB struct {
    *sqlx.DB
}

// Config holds database configuration
type Config struct {
    Host            string
    Port            string
    User            string
    Password        string
    Database        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

// New creates and returns database connection
func New(cfg Config) (*DB, error) {
    // Build DSN (Data Source Name)
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host,
        cfg.Port,
        cfg.User,
        cfg.Password,
        cfg.Database,
        cfg.SSLMode,
    )

    // Open database connection
    sqlDB, err := sqlx.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Set connection pool settings
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &DB{sqlDB}, nil
}

// NewFromEnv creates database connection from environment variables
func NewFromEnv() (*DB, error) {
    cfg := Config{
        Host:            os.Getenv("DB_HOST"),
        Port:            os.Getenv("DB_PORT"),
        User:            os.Getenv("DB_USER"),
        Password:        os.Getenv("DB_PASSWORD"),
        Database:        os.Getenv("DB_NAME"),
        SSLMode:         "disable", // Set dari env atau default
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 5 * time.Minute,
    }

    // Validate required fields
    if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Database == "" {
        return nil, fmt.Errorf("missing required database environment variables")
    }

    return New(cfg)
}

// Close closes database connection
func (db *DB) Close() error {
    if db.DB != nil {
        return db.DB.Close()
    }
    return nil
}

// Health checks database connection status
func (db *DB) Health() error {
    if err := db.Ping(); err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }
    return nil
}

// Stats returns connection pool statistics
func (db *DB) Stats() string {
    stats := db.Stats()
    return fmt.Sprintf(
        "OpenConnections: %d, InUse: %d, Idle: %d, WaitCount: %d",
        stats.OpenConnections,
        stats.InUse,
        stats.Idle,
        stats.WaitCount,
    )
}