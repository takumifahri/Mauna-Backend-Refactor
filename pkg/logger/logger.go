package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// New creates a structured JSON logger configured from LOG_LEVEL.
func New() (*slog.Logger, func() error, error) {
	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	writer := io.Writer(os.Stdout)
	cleanup := func() error { return nil }

	if logFilePath := os.Getenv("LOG_FILE"); logFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(logFilePath), 0o755); err != nil {
			return nil, nil, err
		}

		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, nil, err
		}

		writer = io.MultiWriter(os.Stdout, file)
		cleanup = file.Close
	}

	return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: level,
	})), cleanup, nil
}
