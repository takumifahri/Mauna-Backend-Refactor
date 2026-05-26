package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httphandler "REFACTORING_MAUNA/internal/delivery/http"
	"REFACTORING_MAUNA/pkg/database"
	"REFACTORING_MAUNA/pkg/logger"
	"REFACTORING_MAUNA/pkg/observability"

	"github.com/joho/godotenv"
)

func main() {
	envLoaded := true
	if err := godotenv.Load(); err != nil {
		envLoaded = false
	}

	appLogger, closeLogger, err := logger.New()
	if err != nil {
		slog.Error("failed to initialize logger", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := closeLogger(); err != nil {
			slog.Error("failed to close logger", slog.Any("error", err))
		}
	}()

	slog.SetDefault(appLogger)

	if !envLoaded {
		appLogger.Warn("env file not found")
	}

	shutdownTelemetry, err := observability.InitTracing(context.Background())
	if err != nil {
		appLogger.Error("failed to initialize telemetry", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTelemetry(ctx); err != nil {
			appLogger.Error("failed to shutdown telemetry", slog.Any("error", err))
		}
	}()

	appLogger.Info("telemetry initialized",
		slog.String("service", getEnv("OTEL_SERVICE_NAME", "mauna-backend")),
		slog.String("traces_exporter", getEnv("OTEL_TRACES_EXPORTER", "stdout")),
	)

	port := getEnv("APP_PORT", "8081")

	// Connect to database
	db, err := database.NewFromEnv()
	if err != nil {
		appLogger.Error("failed to connect database", slog.Any("error", err))
		os.Exit(1)
	}

	defer db.Close()

	// Check health
	if err := db.Health(); err != nil {
		appLogger.Error("database health check failed", slog.Any("error", err))
		os.Exit(1)
	}

	appLogger.Info("database connected", slog.String("stats", db.GetStats()))

	// Setup HTTP routes
	mux := http.NewServeMux()
	httphandler.RegisterRoutes(mux, db)

	// Display available routes
	httphandler.PrintRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      observability.HTTPMiddleware(mux, appLogger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("server starting", slog.String("addr", "http://localhost:"+port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	appLogger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("server shutdown error", slog.Any("error", err))
	}

	appLogger.Info("server stopped")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
