package errors

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/go-cli-template/hello-world-cli/internal/logger"
)

// Handler manages error presentation and recovery
type Handler struct {
	Output io.Writer
	Debug  bool
	Color  bool
}

// NewHandler creates a new error handler
func NewHandler() *Handler {
	return &Handler{
		Output: os.Stderr,
		Debug:  false,
		Color:  true,
	}
}

// Handle processes an error and returns the appropriate exit code
func (h *Handler) Handle(err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}

	// Log the full error for debugging
	log := logger.Default()
	log.Error("command failed", "error", err)

	// Present user-friendly error
	h.Present(err)

	// Return appropriate exit code
	return GetExitCode(err)
}

// Present displays an error to the user
func (h *Handler) Present(err error) {
	if err == nil {
		return
	}

	// Build the error message
	var msg strings.Builder

	// Add error icon if color is enabled
	if h.Color {
		msg.WriteString("\033[31m✗\033[0m ")
	} else {
		msg.WriteString("Error: ")
	}

	// Add the main error message
	msg.WriteString(h.getMessage(err))

	// Add suggestions if available
	if suggestion := h.getSuggestion(err); suggestion != "" {
		msg.WriteString("\n\n")
		if h.Color {
			msg.WriteString("\033[33m💡 Suggestion:\033[0m ")
		} else {
			msg.WriteString("Suggestion: ")
		}
		msg.WriteString(suggestion)
	}

	// Add debug information if enabled
	if h.Debug {
		msg.WriteString("\n\n")
		if h.Color {
			msg.WriteString("\033[90mDebug Information:\033[0m\n")
		} else {
			msg.WriteString("Debug Information:\n")
		}
		msg.WriteString(fmt.Sprintf("Error Type: %T\n", err))
		msg.WriteString(fmt.Sprintf("Full Error: %+v\n", err))

		// Add details if it's our custom error
		var appErr *Error
		if errors.As(err, &appErr) && appErr.Details != nil {
			msg.WriteString("Details:\n")
			for k, v := range appErr.Details {
				msg.WriteString(fmt.Sprintf("  %s: %v\n", k, v))
			}
		}
	}

	_, _ = fmt.Fprintln(h.Output, msg.String())
}

// getMessage extracts the user-friendly message from an error
func (h *Handler) getMessage(err error) string {
	// Check for our custom error type first
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr.Message
	}

	// Check specific error types
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		return fmt.Sprintf("Invalid %s: %s", valErr.Field, valErr.Message)
	}

	var cfgErr *ConfigError
	if errors.As(err, &cfgErr) {
		return cfgErr.Error()
	}

	var fileErr *FileError
	if errors.As(err, &fileErr) {
		switch fileErr.Operation {
		case "read":
			return fmt.Sprintf("Cannot read file '%s'", fileErr.Path)
		case "write":
			return fmt.Sprintf("Cannot write to file '%s'", fileErr.Path)
		case "create":
			return fmt.Sprintf("Cannot create file '%s'", fileErr.Path)
		case "delete":
			return fmt.Sprintf("Cannot delete file '%s'", fileErr.Path)
		default:
			return fmt.Sprintf("File operation failed on '%s'", fileErr.Path)
		}
	}

	var netErr *NetworkError
	if errors.As(err, &netErr) {
		if netErr.StatusCode >= 400 {
			return fmt.Sprintf("Request to %s failed with status %d", netErr.URL, netErr.StatusCode)
		}
		return fmt.Sprintf("Network request to %s failed", netErr.URL)
	}

	// Handle common standard errors
	if errors.Is(err, os.ErrNotExist) {
		return "File or directory not found"
	}
	if errors.Is(err, os.ErrPermission) {
		return "Permission denied"
	}

	// Default to the error string
	return err.Error()
}

// getSuggestion returns a helpful suggestion for an error
func (h *Handler) getSuggestion(err error) string {
	var appErr *Error
	if errors.As(err, &appErr) {
		//nolint:exhaustive // We have a default case for unhandled codes
		switch appErr.Code {
		case CodeConfigNotFound:
			return "Run 'hello-world-cli init' to create a configuration file"
		case CodeFilePermission:
			return "Check file permissions or run with appropriate privileges"
		case CodeNetworkTimeout:
			return "Check your internet connection and try again"
		case CodeAuth, CodeUnauthorized:
			return "Run 'hello-world-cli login' to authenticate"
		case CodeMissingArgument:
			return "Use --help to see required arguments"
		default:
			// No specific suggestion for other error codes
		}
	}

	// Validation errors
	if IsValidation(err) {
		return "Check the input format and try again"
	}

	// Network errors
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		if netErr.StatusCode == 404 {
			return "Check the URL and ensure the resource exists"
		}
		if netErr.StatusCode >= 500 {
			return "The server is experiencing issues. Try again later"
		}
	}

	return ""
}

// PanicHandler recovers from panics and converts them to errors
func PanicHandler() {
	if r := recover(); r != nil {
		log := logger.Default()

		// Log the panic with stack trace
		log.Error("panic recovered",
			"panic", r,
			"stack", string(debug.Stack()))

		// Create a user-friendly error
		err := &Error{
			Code:    CodeInternal,
			Message: "An unexpected error occurred",
			Details: map[string]interface{}{
				"panic": fmt.Sprintf("%v", r),
			},
		}

		// Present the error
		handler := NewHandler()
		handler.Present(err)

		// Exit with internal error code
		os.Exit(int(ExitSoftware))
	}
}

// Exit handles error and exits with appropriate code
func Exit(err error) {
	code := defaultHandler.Handle(err)
	os.Exit(int(code))
}

// SetDebug sets the debug mode for error handling
func SetDebug(debugMode bool) {
	defaultHandler.Debug = debugMode
}

var defaultHandler = NewHandler()
