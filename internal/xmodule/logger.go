package xmodule

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger returns a structured logger. Development gets readable text;
// production gets JSON for log collectors.
func NewLogger(level string) *slog.Logger {
	var parsed slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		parsed = slog.LevelDebug
	case "warn":
		parsed = slog.LevelWarn
	case "error":
		parsed = slog.LevelError
	default:
		parsed = slog.LevelInfo
	}
	options := &slog.HandlerOptions{Level: parsed}
	if strings.EqualFold(os.Getenv("APP_ENV"), "production") || os.Getenv("APP_ENV") == "" && os.Getenv("LISTEN") != "" && !strings.HasPrefix(os.Getenv("LISTEN"), "127.") {
		return slog.New(slog.NewJSONHandler(os.Stdout, options))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, options))
}
