package errors

import (
	"errors"
	"fmt"
)

// Error represents a structured application error with context
type Error struct {
	Code    ErrorCode              // Machine-readable error code
	Message string                 // User-friendly message
	Err     error                  // Wrapped error
	Details map[string]interface{} // Additional context
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.Err
}

// Is implements errors.Is
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// WithDetails adds or updates details on the error
func (e *Error) WithDetails(key string, value interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// New creates a new application error
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrapf wraps an error with a formatted message
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Err:     err,
	}
}

// ValidationError represents an input validation error
type ValidationError struct {
	Field   string      // Field that failed validation
	Value   interface{} // The invalid value
	Message string      // Validation message
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// ConfigError represents a configuration error
type ConfigError struct {
	Key     string // Configuration key
	Value   string // Problematic value
	Message string // Error message
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error for %s: %s", e.Key, e.Message)
}

// FileError represents a file operation error
type FileError struct {
	Path      string // File path
	Operation string // Operation that failed (read, write, delete, etc.)
	Err       error  // Underlying error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file %s failed for %s: %v", e.Operation, e.Path, e.Err)
}

func (e *FileError) Unwrap() error {
	return e.Err
}

// NetworkError represents a network-related error
type NetworkError struct {
	URL        string // URL or address
	Operation  string // GET, POST, etc.
	StatusCode int    // HTTP status code if applicable
	Err        error  // Underlying error
}

func (e *NetworkError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("network error %s %s: status %d", e.Operation, e.URL, e.StatusCode)
	}
	return fmt.Sprintf("network error %s %s: %v", e.Operation, e.URL, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// Common error checking helpers

// IsValidation checks if an error is a validation error
func IsValidation(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// IsConfig checks if an error is a configuration error
func IsConfig(err error) bool {
	var ce *ConfigError
	return errors.As(err, &ce)
}

// IsFile checks if an error is a file error
func IsFile(err error) bool {
	var fe *FileError
	return errors.As(err, &fe)
}

// IsNetwork checks if an error is a network error
func IsNetwork(err error) bool {
	var ne *NetworkError
	return errors.As(err, &ne)
}

// IsCode checks if an error has a specific error code
func IsCode(err error, code ErrorCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}
