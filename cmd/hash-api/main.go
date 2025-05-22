package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/derticom/hash-api/internal/config"
)

func main() {
	cfg := config.NewConfig()

	log, err := setupLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to setup logger: %+v", err))
	}

	log.Info("Logger initialized")

	// TODO: init storage

	// TODO: init router

	// TODO: init server
}

func setupLogger(cfg *config.Config) (*slog.Logger, error) {
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			},
		),
	)

	return logger, nil
}
