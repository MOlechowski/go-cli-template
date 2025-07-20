package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// Logger is the interface that wraps basic logging methods.
// This allows for easy swapping of logging implementations.
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)

	With(fields ...any) Logger
	WithError(err error) Logger
}

// Config holds logger configuration
type Config struct {
	Level   string // debug, info, warn, error
	Format  string // text, json
	Output  string // stdout, stderr
	NoColor bool   // disable color in text format
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:   "info",
		Format:  "text",
		Output:  "stderr",
		NoColor: false,
	}
}

// slogLogger wraps slog.Logger to implement our Logger interface
type slogLogger struct {
	logger *slog.Logger
}

// New creates a new logger with the given configuration
func New(cfg Config) Logger {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var output *os.File
	switch cfg.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output = os.Stderr
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	case "text":
		if cfg.NoColor || !isTerminal(output) {
			handler = slog.NewTextHandler(output, opts)
		} else {
			handler = NewColorHandler(output, opts)
		}
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	return &slogLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a message at debug level
func (l *slogLogger) Debug(msg string, fields ...any) {
	l.logger.Debug(msg, fields...)
}

// Info logs a message at info level
func (l *slogLogger) Info(msg string, fields ...any) {
	l.logger.Info(msg, fields...)
}

// Warn logs a message at warn level
func (l *slogLogger) Warn(msg string, fields ...any) {
	l.logger.Warn(msg, fields...)
}

// Error logs a message at error level
func (l *slogLogger) Error(msg string, fields ...any) {
	l.logger.Error(msg, fields...)
}

// With returns a new logger with additional fields
func (l *slogLogger) With(fields ...any) Logger {
	return &slogLogger{
		logger: l.logger.With(fields...),
	}
}

// WithError returns a new logger with an error field
func (l *slogLogger) WithError(err error) Logger {
	return l.With("error", err)
}

// FromContext returns the logger from context, or a default logger if not found
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return logger
	}
	return New(DefaultConfig())
}

// WithContext returns a new context with the logger attached
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

type loggerKey struct{}

// isTerminal checks if the file descriptor is a terminal
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	// Simple check - in production you might want to use golang.org/x/term
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// Global logger instance for convenience
var defaultLogger = New(DefaultConfig())

// SetDefault sets the default global logger
func SetDefault(l Logger) {
	defaultLogger = l
}

// Default returns the default global logger
func Default() Logger {
	return defaultLogger
}

// Debug logs a message at debug level using the default logger
func Debug(msg string, fields ...any) { defaultLogger.Debug(msg, fields...) }

// Info logs a message at info level using the default logger
func Info(msg string, fields ...any) { defaultLogger.Info(msg, fields...) }

// Warn logs a message at warn level using the default logger
func Warn(msg string, fields ...any) { defaultLogger.Warn(msg, fields...) }

// Error logs a message at error level using the default logger
func Error(msg string, fields ...any) { defaultLogger.Error(msg, fields...) }
