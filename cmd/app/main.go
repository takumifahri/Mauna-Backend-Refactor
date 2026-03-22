package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "REFACTORING_MAUNA/pkg/database"
    httphandler "REFACTORING_MAUNA/internal/delivery/http"

    "github.com/joho/godotenv"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  Warning: .env file not found")
    }

    // Connect to database
    db, err := database.NewFromEnv()
    if err != nil {
        log.Fatalf("❌ Failed to connect database: %v", err)
    }
    defer db.Close()

    // Check health
    if err := db.Health(); err != nil {
        log.Fatalf("❌ Database health check failed: %v", err)
    }

    log.Println("✅ Database connected successfully!")
    log.Println("📊 " + db.GetStats())

    // Setup HTTP routes
    mux := http.NewServeMux()
    httphandler.RegisterRoutes(mux, db)

    // Display available routes
    httphandler.PrintRoutes()

    // Create HTTP server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in goroutine
    go func() {
        log.Printf("🚀 Server starting on http://localhost:8080\n")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("❌ Server error: %v", err)
        }
    }()

    // Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    <-sigChan
    log.Println("\n⏸️  Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("❌ Server shutdown error: %v", err)
    }

    log.Println("✅ Server stopped")
}