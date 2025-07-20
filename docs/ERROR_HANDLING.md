# Error Handling Guide

This CLI template includes a comprehensive error handling system that provides:
- Structured error types with context
- User-friendly error messages
- Proper exit codes for scripting
- Panic recovery
- Debug mode for troubleshooting

## Quick Start

### Using Custom Error Types

```go
import "github.com/go-cli-template/hello-world-cli/internal/errors"

// Return a validation error
if name == "" {
    return &errors.ValidationError{
        Field:   "name",
        Value:   name,
        Message: "name cannot be empty",
    }
}

// Wrap an existing error with context
data, err := os.ReadFile(path)
if err != nil {
    return errors.Wrap(err, errors.CodeFileRead, "failed to read configuration")
}

// Create a new error with details
return errors.New(errors.CodeConfigInvalid, "invalid configuration").
    WithDetails("file", configPath).
    WithDetails("error", "missing required field")
```

## Error Types

### Base Error Type

The `Error` struct provides structured errors with:
- `Code`: Machine-readable error code
- `Message`: User-friendly message
- `Err`: Wrapped underlying error
- `Details`: Additional context as key-value pairs

### Specialized Error Types

1. **ValidationError**: For input validation failures
   ```go
   &errors.ValidationError{
       Field:   "email",
       Value:   userInput,
       Message: "invalid email format",
   }
   ```

2. **ConfigError**: For configuration issues
   ```go
   &errors.ConfigError{
       Key:     "database.url",
       Value:   dbURL,
       Message: "invalid connection string",
   }
   ```

3. **FileError**: For file operations
   ```go
   &errors.FileError{
       Path:      "/path/to/file",
       Operation: "read",
       Err:       originalErr,
   }
   ```

4. **NetworkError**: For network operations
   ```go
   &errors.NetworkError{
       URL:        "https://api.example.com",
       Operation:  "GET",
       StatusCode: 404,
       Err:        originalErr,
   }
   ```

## Error Codes

Common error codes and their meanings:

```go
// General
CodeUnknown         // Unknown error
CodeInternal        // Internal software error
CodeTimeout         // Operation timed out
CodeCanceled        // Operation was canceled

// Validation
CodeInvalidInput    // Invalid input provided
CodeMissingArgument // Required argument missing
CodeValidation      // General validation error

// Configuration
CodeConfigNotFound  // Config file not found
CodeConfigInvalid   // Invalid configuration

// File I/O
CodeFileNotFound    // File doesn't exist
CodeFilePermission  // Permission denied
CodeFileRead        // Read operation failed

// Network
CodeNetworkTimeout  // Network timeout
CodeNetworkConnect  // Connection failed
```

## Exit Codes

The error handler automatically maps errors to appropriate exit codes:

| Exit Code | Meaning | Used For |
|-----------|---------|----------|
| 0 | Success | Successful completion |
| 1 | General Error | Default for unknown errors |
| 2 | Misuse | Command line usage error |
| 65 | Data Error | Invalid data format |
| 66 | No Input | Cannot open input |
| 69 | Unavailable | Service unavailable |
| 70 | Software | Internal software error |
| 74 | I/O Error | Input/output error |
| 77 | No Permission | Permission denied |
| 78 | Config | Configuration error |

## Error Presentation

### Default Mode

Users see clean, helpful error messages:

```
✗ Invalid name: name cannot be empty

💡 Suggestion: Use --help to see required arguments
```

### Debug Mode

With `--debug` flag, users see detailed error information:

```
✗ Invalid name: name cannot be empty

💡 Suggestion: Use --help to see required arguments

Debug Information:
Error Type: *errors.ValidationError
Full Error: validation failed for name: name cannot be empty
Details:
  field: name
  value:
```

## Panic Recovery

The application automatically recovers from panics:

```go
func main() {
    // Panic recovery is set up automatically
    defer errors.PanicHandler()
    
    // Your code here
}
```

If a panic occurs, users see:
```
✗ An unexpected error occurred

💡 Suggestion: Please report this issue
```

While the full panic details are logged for debugging.

## Best Practices

### 1. Use Specific Error Types

```go
// Good - specific error type with context
return &errors.FileError{
    Path:      configPath,
    Operation: "read",
    Err:       err,
}

// Avoid - generic error without context
return fmt.Errorf("failed to read file")
```

### 2. Add Helpful Context

```go
// Good - includes what was being done
return errors.Wrap(err, errors.CodeNetwork, "failed to fetch user data from API")

// Better - includes even more context
return errors.Wrapf(err, errors.CodeNetwork, 
    "failed to fetch user %s from API endpoint %s", userID, endpoint)
```

### 3. Provide User-Friendly Messages

```go
// Good - tells user what's wrong and what to do
return errors.New(errors.CodeConfigNotFound, 
    "Configuration file not found. Run 'myapp init' to create one")

// Avoid - technical jargon
return errors.New(errors.CodeConfigNotFound, 
    "ENOENT: no such file or directory")
```

### 4. Use Error Checking Helpers

```go
// Check error type
if errors.IsValidation(err) {
    // Handle validation error specifically
}

// Check error code
if errors.IsCode(err, errors.CodeTimeout) {
    // Retry the operation
}
```

## Examples

### Command Implementation

```go
func runCommand(cmd *cobra.Command, opts *Options) error {
    // Validate input
    if opts.Input == "" {
        return &errors.ValidationError{
            Field:   "input",
            Value:   opts.Input,
            Message: "input file is required",
        }
    }
    
    // Read file
    data, err := os.ReadFile(opts.Input)
    if err != nil {
        return &errors.FileError{
            Path:      opts.Input,
            Operation: "read",
            Err:       err,
        }
    }
    
    // Process data
    result, err := processData(data)
    if err != nil {
        return errors.Wrap(err, errors.CodeInternal, 
            "failed to process input data")
    }
    
    // Success
    return nil
}
```

### Main Function

```go
func main() {
    // Set up panic recovery
    defer errors.PanicHandler()
    
    // Execute command
    if err := cli.Execute(); err != nil {
        // Error handler shows user-friendly message
        // and exits with appropriate code
        errors.Exit(err)
    }
}
```

## Testing Errors

```go
func TestCommandError(t *testing.T) {
    // Test validation error
    err := runCommand(cmd, &Options{Input: ""})
    
    // Check it's a validation error
    if !errors.IsValidation(err) {
        t.Errorf("expected validation error, got %T", err)
    }
    
    // Check exit code
    code := errors.GetExitCode(err)
    if code != errors.ExitDataError {
        t.Errorf("expected exit code %d, got %d", 
            errors.ExitDataError, code)
    }
}
```

## Migration Guide

To migrate existing error handling:

1. Replace `fmt.Errorf` with appropriate error types
2. Use `errors.Wrap` instead of `fmt.Errorf("...: %w", err)`
3. Add error codes for different failure modes
4. Ensure main.go uses `errors.Exit(err)`

Before:
```go
if name == "" {
    return fmt.Errorf("name is required")
}
```

After:
```go
if name == "" {
    return &errors.ValidationError{
        Field:   "name",
        Value:   name,
        Message: "name is required",
    }
}
```