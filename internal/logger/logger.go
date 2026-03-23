package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	once sync.Once
)

func NewLogger(level string) {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	case "error":
		logLevel = slog.LevelError
	case "warning":
		logLevel = slog.LevelWarn
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
