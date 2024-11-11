// Package logger provides a simple logging interface that complies with the partial method set of slog.Logger.
// It also defines a standard context key name for logger instances.
package logger

import (
	"context"
	"log/slog"
	"os"
)

type ContextKey string

var LoggerContextKey = ContextKey("logger")

// Logger is a simple logging interface that complies with the partial method set of slog.Logger.
type Logger interface {
	Enabled(ctx context.Context, level slog.Level) bool // Enabled returns true if the given log level is enabled.
	Debug(msg string, fields ...any)                    // Debug logs a message at the debug level.
	Error(msg string, fields ...any)                    // Error logs a message at the error level.
	Info(msg string, fields ...any)                     // Info logs a message at the info level.
	Warn(msg string, fields ...any)                     // Warn logs a message at the warn level.
}

// GetLogger returns the logger from the context or a default logger if none is found.
func GetLogger(ctx context.Context) Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(Logger); ok {
		return logger
	}
	return NewDefaultLogger()
}

// NewDefaultLogger returns a new error logger that writes to stdout.
func NewDefaultLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

// NewDebugLogger returns a new debug logger with source that writes to stdout.
func NewDebugLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
}
