package logger

import (
	"log/slog"
	"os"
)

// Create new logger instance and return it. Will be bound to the server in server.go
func New() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
}
