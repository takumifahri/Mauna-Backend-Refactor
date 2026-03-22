package database

import (
    "context"
    "fmt"
    "math"
    "os"
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
    MaxRetries      int           // Tambah
    RetryDelay      time.Duration // Tambah
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

    var sqlDB *sqlx.DB
    var err error

    // Retry logic dengan exponential backoff
    maxRetries := cfg.MaxRetries
    if maxRetries == 0 {
        maxRetries = 3
    }

    for attempt := 0; attempt < maxRetries; attempt++ {
        sqlDB, err = sqlx.Open("postgres", dsn)
        if err == nil {
            // Test connection
            if pingErr := sqlDB.Ping(); pingErr == nil {
                break
            } else {
                err = pingErr
            }
        }

        if attempt < maxRetries-1 {
            backoff := time.Duration(math.Pow(2, float64(attempt))) * (cfg.RetryDelay)
            fmt.Printf("⚠️  Database connection failed (attempt %d/%d), retrying in %v...\n", attempt+1, maxRetries, backoff)
            time.Sleep(backoff)
        }
    }

    if err != nil {
        return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
    }

    // Set connection pool settings
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

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
        SSLMode:         "disable",
        MaxOpenConns:    25,
        MaxIdleConns:    5,
        ConnMaxLifetime: 5 * time.Minute,
        MaxRetries:      3,
        RetryDelay:      1 * time.Second,
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
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }
    return nil
}

// GetStats returns connection pool statistics
func (db *DB) GetStats() string {
    stats := db.DB.Stats()
    return fmt.Sprintf(
        "OpenConnections: %d, InUse: %d, Idle: %d, WaitCount: %d",
        stats.OpenConnections,
        stats.InUse,
        stats.Idle,
        stats.WaitCount,
    )
}

// WithTx executes a function within a transaction
// Automatically handles rollback on error, commit on success
func (db *DB) WithTx(fn func(*sqlx.Tx) error) error {
    tx, err := db.Beginx()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("transaction failed with error: %w, rollback also failed: %w", err, rbErr)
        }
        return fmt.Errorf("transaction failed: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}

// WithTxContext executes a function within a transaction with context
func (db *DB) WithTxContext(ctx context.Context, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }

    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("transaction failed with error: %w, rollback also failed: %w", err, rbErr)
        }
        return fmt.Errorf("transaction failed: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}