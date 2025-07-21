package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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
	Level     string // debug, info, warn, error
	Format    string // text, json
	Output    string // stdout, stderr
	NoColor   bool   // disable color in text format
	AddSource bool   // include file:line in output (auto-enabled for debug level)
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
	level := parseLogLevel(cfg.Level, &cfg)
	opts := createHandlerOptions(level, cfg)
	output := getOutput(cfg.Output)
	handler := createHandler(cfg, output, opts)

	return &slogLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a message at debug level
func (l *slogLogger) Debug(msg string, fields ...any) {
	// Skip 2 frames: this method and the caller's location
	l.logWithCaller(slog.LevelDebug, msg, fields...)
}

// Info logs a message at info level
func (l *slogLogger) Info(msg string, fields ...any) {
	l.logWithCaller(slog.LevelInfo, msg, fields...)
}

// Warn logs a message at warn level
func (l *slogLogger) Warn(msg string, fields ...any) {
	l.logWithCaller(slog.LevelWarn, msg, fields...)
}

// Error logs a message at error level
func (l *slogLogger) Error(msg string, fields ...any) {
	l.logWithCaller(slog.LevelError, msg, fields...)
}

// logWithCaller logs with the correct caller information
func (l *slogLogger) logWithCaller(level slog.Level, msg string, fields ...any) {
	var pcs [1]uintptr
	// Skip 3 frames to get the real caller:
	// 1. runtime.Callers
	// 2. this method (logWithCaller)
	// 3. the wrapper method (Debug, Info, etc.)
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(fields...)
	_ = l.logger.Handler().Handle(context.Background(), r)
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

// parseLogLevel converts string level to slog.Level and handles debug special case
func parseLogLevel(levelStr string, cfg *Config) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		// Auto-enable source location for debug level
		cfg.AddSource = true
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// createHandlerOptions creates slog.HandlerOptions with source replacement
func createHandlerOptions(level slog.Level, cfg Config) *slog.HandlerOptions {
	// Get current working directory for relative paths
	wd, _ := os.Getwd()

	// ReplaceAttr to show relative paths instead of absolute
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey && cfg.AddSource {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				// Show relative path from working directory
				if wd != "" {
					if relPath, err := filepath.Rel(wd, source.File); err == nil {
						source.File = relPath
					}
				}
				// Remove function name to keep output concise
				source.Function = ""
			}
		}
		return a
	}

	return &slog.HandlerOptions{
		Level:       level,
		AddSource:   cfg.AddSource,
		ReplaceAttr: replaceAttr,
	}
}

// getOutput returns the appropriate output file
func getOutput(outputStr string) *os.File {
	switch outputStr {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return os.Stderr
	}
}

// createHandler creates the appropriate slog.Handler based on configuration
func createHandler(cfg Config, output *os.File, opts *slog.HandlerOptions) slog.Handler {
	switch strings.ToLower(cfg.Format) {
	case "json":
		return slog.NewJSONHandler(output, opts)
	case "text":
		if cfg.NoColor || !isTerminal(output) {
			return slog.NewTextHandler(output, opts)
		}
		return NewColorHandler(output, opts)
	default:
		return slog.NewTextHandler(output, opts)
	}
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
func Debug(msg string, fields ...any) {
	if sl, ok := defaultLogger.(*slogLogger); ok {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // Skip runtime.Callers and this function
		r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0])
		r.Add(fields...)
		_ = sl.logger.Handler().Handle(context.Background(), r)
	} else {
		defaultLogger.Debug(msg, fields...)
	}
}

// Info logs a message at info level using the default logger
func Info(msg string, fields ...any) {
	if sl, ok := defaultLogger.(*slogLogger); ok {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:])
		r := slog.NewRecord(time.Now(), slog.LevelInfo, msg, pcs[0])
		r.Add(fields...)
		_ = sl.logger.Handler().Handle(context.Background(), r)
	} else {
		defaultLogger.Info(msg, fields...)
	}
}

// Warn logs a message at warn level using the default logger
func Warn(msg string, fields ...any) {
	if sl, ok := defaultLogger.(*slogLogger); ok {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:])
		r := slog.NewRecord(time.Now(), slog.LevelWarn, msg, pcs[0])
		r.Add(fields...)
		_ = sl.logger.Handler().Handle(context.Background(), r)
	} else {
		defaultLogger.Warn(msg, fields...)
	}
}

// Error logs a message at error level using the default logger
func Error(msg string, fields ...any) {
	if sl, ok := defaultLogger.(*slogLogger); ok {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:])
		r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
		r.Add(fields...)
		_ = sl.logger.Handler().Handle(context.Background(), r)
	} else {
		defaultLogger.Error(msg, fields...)
	}
}
