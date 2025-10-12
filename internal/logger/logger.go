package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a new structured logger using slog.
func New(level, serviceName, env string) *slog.Logger {
	var loglevel slog.Level

	// define log level based on input string
	switch strings.ToUpper(level) {
	case "DEBUG":
		loglevel = slog.LevelDebug
	case "INFO":
		loglevel = slog.LevelInfo
	case "WARN":
		loglevel = slog.LevelWarn
	case "ERROR":
		loglevel = slog.LevelError
	default:
		loglevel = slog.LevelInfo
	}

	var handler slog.Handler

	// define handler based on environment
	if env == "development" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     loglevel,
			AddSource: true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: loglevel,
		})
	}

	// add default fields
	logger := slog.New(handler).With(
		slog.String("service", serviceName),
		slog.String("env", env),
	)

	return logger
}
