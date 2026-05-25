package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

func Setup() error {
	logDir := "logs"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "app.log")

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(
		slog.NewMultiHandler(
			consoleHandler,
			fileHandler,
		),
	)

	slog.SetDefault(logger)

	return nil
}
