package errors

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestHandler_Present(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		debug     bool
		wantInOut []string // strings that should be in output
	}{
		{
			name: "simple error message",
			err:  New(CodeInvalidInput, "Invalid input provided"),
			wantInOut: []string{
				"Invalid input provided",
			},
		},
		{
			name: "validation error",
			err: &ValidationError{
				Field:   "email",
				Value:   "invalid",
				Message: "must be a valid email",
			},
			wantInOut: []string{
				"Invalid email: must be a valid email",
			},
		},
		{
			name: "error with suggestion",
			err:  New(CodeConfigNotFound, "Configuration file not found"),
			wantInOut: []string{
				"Configuration file not found",
				"Suggestion:",
				"Run 'hello-world-cli init'",
			},
		},
		{
			name:  "debug mode shows details",
			err:   New(CodeInternal, "Internal error").WithDetails("component", "database"),
			debug: true,
			wantInOut: []string{
				"Internal error",
				"Debug Information:",
				"Error Type:",
				"component: database",
			},
		},
		{
			name: "file error with nice message",
			err: &FileError{
				Path:      "/etc/config.yaml",
				Operation: "read",
			},
			wantInOut: []string{
				"Cannot read file '/etc/config.yaml'",
			},
		},
		{
			name: "network error with status code",
			err: &NetworkError{
				URL:        "https://api.example.com",
				Operation:  "GET",
				StatusCode: 404,
			},
			wantInOut: []string{
				"Request to https://api.example.com failed with status 404",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := &Handler{
				Output: &buf,
				Debug:  tt.debug,
				Color:  false, // Disable color for testing
			}

			h.Present(tt.err)

			output := buf.String()
			for _, want := range tt.wantInOut {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestHandler_GetMessage(t *testing.T) {
	h := &Handler{}

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "custom error message",
			err:  New(CodeInternal, "Custom message"),
			want: "Custom message",
		},
		{
			name: "validation error",
			err:  &ValidationError{Field: "age", Message: "must be positive"},
			want: "Invalid age: must be positive",
		},
		{
			name: "os.ErrNotExist",
			err:  os.ErrNotExist,
			want: "File or directory not found",
		},
		{
			name: "os.ErrPermission",
			err:  os.ErrPermission,
			want: "Permission denied",
		},
		{
			name: "generic error",
			err:  fmt.Errorf("something went wrong"),
			want: "something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.getMessage(tt.err); got != tt.want {
				t.Errorf("getMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GetSuggestion(t *testing.T) {
	h := &Handler{}

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "config not found",
			err:  New(CodeConfigNotFound, "Config missing"),
			want: "Run 'hello-world-cli init' to create a configuration file",
		},
		{
			name: "permission error",
			err:  New(CodeFilePermission, "Permission denied"),
			want: "Check file permissions or run with appropriate privileges",
		},
		{
			name: "network timeout",
			err:  New(CodeNetworkTimeout, "Request timed out"),
			want: "Check your internet connection and try again",
		},
		{
			name: "auth error",
			err:  New(CodeUnauthorized, "Not authenticated"),
			want: "Run 'hello-world-cli login' to authenticate",
		},
		{
			name: "validation error",
			err:  &ValidationError{Field: "test", Message: "invalid"},
			want: "Check the input format and try again",
		},
		{
			name: "network 404",
			err:  &NetworkError{StatusCode: 404},
			want: "Check the URL and ensure the resource exists",
		},
		{
			name: "network 500",
			err:  &NetworkError{StatusCode: 500},
			want: "The server is experiencing issues. Try again later",
		},
		{
			name: "no suggestion",
			err:  fmt.Errorf("unknown error"),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.getSuggestion(tt.err); got != tt.want {
				t.Errorf("getSuggestion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		wantExitCode ExitCode
	}{
		{
			name:         "nil error returns success",
			err:          nil,
			wantExitCode: ExitSuccess,
		},
		{
			name:         "validation error returns data error exit",
			err:          &ValidationError{Field: "test", Message: "invalid"},
			wantExitCode: ExitDataError,
		},
		{
			name:         "config error returns config exit",
			err:          New(CodeConfig, "Config error"),
			wantExitCode: ExitConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := &Handler{
				Output: &buf,
				Debug:  false,
				Color:  false,
			}

			if got := h.Handle(tt.err); got != tt.wantExitCode {
				t.Errorf("Handle() = %v, want %v", got, tt.wantExitCode)
			}
		})
	}
}
