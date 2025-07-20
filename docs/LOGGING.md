# Logging Guide

This CLI template includes a flexible logging system based on Go's standard `log/slog` package, with support for different log levels, formats, and outputs.

## Quick Start

The logger is automatically configured and available in all commands:

```go
func runMyCommand(cmd *cobra.Command, opts *Options) error {
    log := logger.FromContext(cmd.Context())
    
    log.Debug("command started", "option", opts.Value)
    log.Info("processing request", "user", opts.User)
    log.Warn("deprecated feature used")
    log.Error("operation failed", "error", err)
    
    return nil
}
```

## Configuration

The logger can be configured through multiple methods (in order of precedence):

1. **Environment Variables** (highest priority)
   ```bash
   LOG_LEVEL=debug LOG_FORMAT=json ./hello-world-cli
   ```

2. **Command Line Flags**
   ```bash
   ./hello-world-cli --debug                    # Sets log level to debug
   ./hello-world-cli --log-level=warn          # Sets specific log level
   ./hello-world-cli --log-format=json         # JSON output format
   ```

3. **Configuration File** (`.hello-world-cli.yaml`)
   ```yaml
   log:
     level: debug
     format: json
   ```

## Log Levels

Available log levels (from most to least verbose):
- `debug` - Detailed information for debugging
- `info` - General informational messages (default)
- `warn` - Warning messages
- `error` - Error messages only

## Output Formats

### Text Format (Default)
Colorized output for development with easy-to-read formatting:
```
15:04:05 INFO  Processing user request name="Alice" language="en"
15:04:05 DEBUG Generated greeting message="Hello, Alice!"
15:04:05 ERROR Failed to save result error="permission denied"
```

### JSON Format
Structured output for production and log aggregation:
```json
{"time":"2024-01-20T15:04:05Z","level":"INFO","msg":"Processing user request","name":"Alice","language":"en"}
{"time":"2024-01-20T15:04:05Z","level":"ERROR","msg":"Failed to save result","error":"permission denied"}
```

## Using the Logger

### Basic Usage

```go
// Get logger from context (recommended)
log := logger.FromContext(ctx)

// Or use the global logger
logger.Info("global log message")
```

### Adding Context Fields

```go
// Add fields to a specific log entry
log.Info("user action", 
    "user_id", 123,
    "action", "login",
    "ip", "192.168.1.1",
)

// Create a logger with persistent fields
userLog := log.With("user_id", 123, "session", "abc123")
userLog.Info("logged in")           // Includes user_id and session
userLog.Info("viewed dashboard")    // Includes user_id and session
```

### Error Logging

```go
// Log with error
if err != nil {
    log.Error("operation failed", "error", err)
    // or
    log.WithError(err).Error("operation failed")
}
```

## Best Practices

1. **Use appropriate log levels**:
   - `Debug`: Detailed execution flow, variable values
   - `Info`: Important business events, state changes
   - `Warn`: Recoverable issues, deprecated features
   - `Error`: Failures that need attention

2. **Add context to logs**:
   ```go
   log.Info("processing payment",
       "user_id", userID,
       "amount", amount,
       "currency", currency,
   )
   ```

3. **Use structured fields** instead of string formatting:
   ```go
   // Good
   log.Info("user logged in", "user_id", userID, "ip", ipAddr)
   
   // Avoid
   log.Info(fmt.Sprintf("User %d logged in from %s", userID, ipAddr))
   ```

4. **Create scoped loggers** for different components:
   ```go
   dbLog := log.With("component", "database")
   apiLog := log.With("component", "api")
   ```

## Switching to Alternative Loggers

While `slog` is the default, you can easily switch to other popular loggers:

### Using Zap

1. Install zap:
   ```bash
   go get go.uber.org/zap
   ```

2. Create a zap adapter in `internal/logger/zap_adapter.go`:
   ```go
   type zapLogger struct {
       logger *zap.SugaredLogger
   }
   
   func NewZapLogger(cfg Config) Logger {
       // Configure zap based on cfg
       // Return zapLogger that implements Logger interface
   }
   ```

3. Update logger initialization in your code.

### Using Zerolog

1. Install zerolog:
   ```bash
   go get github.com/rs/zerolog
   ```

2. Create a zerolog adapter following the same pattern.

## Performance Considerations

- The default `slog` logger is efficient for most CLI applications
- For high-performance scenarios, consider:
  - Using `zerolog` for zero-allocation logging
  - Using `zap` for high-throughput applications
  - Disabling debug logging in production
  - Using JSON format for better parsing performance

## Troubleshooting

### No log output?
- Check the log level - debug messages won't show at info level
- Ensure you're looking at stderr (default output)
- Verify logger initialization in PersistentPreRunE

### Colors not showing?
- Colors are automatically disabled when output is redirected
- Use `--log-format=text` to force text format
- Check if your terminal supports ANSI colors

### Too much/little logging?
- Adjust log level: `--log-level=warn` for less output
- Use `--debug` flag for maximum verbosity
- Set `LOG_LEVEL` environment variable for persistent setting