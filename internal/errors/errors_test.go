package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		wantMsg  string
		wantCode ErrorCode
	}{
		{
			name: "simple error",
			err: &Error{
				Code:    CodeInvalidInput,
				Message: "invalid input provided",
			},
			wantMsg:  "invalid input provided",
			wantCode: CodeInvalidInput,
		},
		{
			name: "error with wrapped error",
			err: &Error{
				Code:    CodeFileRead,
				Message: "cannot read file",
				Err:     fmt.Errorf("permission denied"),
			},
			wantMsg:  "cannot read file: permission denied",
			wantCode: CodeFileRead,
		},
		{
			name: "error with details",
			err: func() *Error {
				e := &Error{
					Code:    CodeConfig,
					Message: "configuration error",
				}
				_ = e.WithDetails("key", "database.url").WithDetails("value", "invalid")
				return e
			}(),
			wantMsg:  "configuration error",
			wantCode: CodeConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("Error() = %v, want %v", got, tt.wantMsg)
			}
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %v, want %v", tt.err.Code, tt.wantCode)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	err := &Error{
		Code:    CodeInternal,
		Message: "wrapped error",
		Err:     baseErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestErrorIs(t *testing.T) {
	err1 := &Error{Code: CodeInvalidInput}
	err2 := &Error{Code: CodeInvalidInput}
	err3 := &Error{Code: CodeFileRead}

	if !err1.Is(err2) {
		t.Error("Expected errors with same code to match")
	}
	if err1.Is(err3) {
		t.Error("Expected errors with different codes not to match")
	}
	if err1.Is(fmt.Errorf("other error")) {
		t.Error("Expected error not to match non-Error type")
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Value:   "invalid",
		Message: "invalid email format",
	}

	want := "validation failed for email: invalid email format"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}
}

func TestConfigError(t *testing.T) {
	err := &ConfigError{
		Key:     "database.host",
		Value:   "invalid:port",
		Message: "invalid host format",
	}

	want := "config error for database.host: invalid host format"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}
}

func TestFileError(t *testing.T) {
	baseErr := fmt.Errorf("permission denied")
	err := &FileError{
		Path:      "/etc/config.yaml",
		Operation: "read",
		Err:       baseErr,
	}

	want := "file read failed for /etc/config.yaml: permission denied"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}

	if unwrapped := err.Unwrap(); unwrapped != baseErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestNetworkError(t *testing.T) {
	tests := []struct {
		name string
		err  *NetworkError
		want string
	}{
		{
			name: "with status code",
			err: &NetworkError{
				URL:        "https://api.example.com",
				Operation:  "GET",
				StatusCode: 404,
			},
			want: "network error GET https://api.example.com: status 404",
		},
		{
			name: "with underlying error",
			err: &NetworkError{
				URL:       "https://api.example.com",
				Operation: "POST",
				Err:       fmt.Errorf("connection refused"),
			},
			want: "network error POST https://api.example.com: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorTypeCheckers(t *testing.T) {
	valErr := &ValidationError{Field: "test", Message: "test"}
	cfgErr := &ConfigError{Key: "test", Message: "test"}
	fileErr := &FileError{Path: "test", Operation: "test"}
	netErr := &NetworkError{URL: "test", Operation: "test"}

	// Test IsValidation
	if !IsValidation(valErr) {
		t.Error("IsValidation should return true for ValidationError")
	}
	if IsValidation(cfgErr) {
		t.Error("IsValidation should return false for non-ValidationError")
	}

	// Test IsConfig
	if !IsConfig(cfgErr) {
		t.Error("IsConfig should return true for ConfigError")
	}
	if IsConfig(valErr) {
		t.Error("IsConfig should return false for non-ConfigError")
	}

	// Test IsFile
	if !IsFile(fileErr) {
		t.Error("IsFile should return true for FileError")
	}
	if IsFile(netErr) {
		t.Error("IsFile should return false for non-FileError")
	}

	// Test IsNetwork
	if !IsNetwork(netErr) {
		t.Error("IsNetwork should return true for NetworkError")
	}
	if IsNetwork(fileErr) {
		t.Error("IsNetwork should return false for non-NetworkError")
	}
}

func TestIsCode(t *testing.T) {
	err := &Error{Code: CodeTimeout}
	wrappedErr := fmt.Errorf("wrapped: %w", err)

	if !IsCode(err, CodeTimeout) {
		t.Error("IsCode should return true for matching code")
	}
	if IsCode(err, CodeInternal) {
		t.Error("IsCode should return false for non-matching code")
	}
	if !IsCode(wrappedErr, CodeTimeout) {
		t.Error("IsCode should work with wrapped errors")
	}
	if IsCode(fmt.Errorf("other error"), CodeTimeout) {
		t.Error("IsCode should return false for non-Error types")
	}
}

func TestWrapFunctions(t *testing.T) {
	baseErr := fmt.Errorf("base error")

	// Test Wrap
	wrapped := Wrap(baseErr, CodeInternal, "operation failed")
	if wrapped.Code != CodeInternal {
		t.Errorf("Code = %v, want %v", wrapped.Code, CodeInternal)
	}
	if wrapped.Message != "operation failed" {
		t.Errorf("Message = %v, want %v", wrapped.Message, "operation failed")
	}
	if !errors.Is(wrapped, baseErr) {
		t.Error("Wrapped error should contain base error")
	}

	// Test Wrapf
	wrappedf := Wrapf(baseErr, CodeFileRead, "failed to read %s", "config.yaml")
	if wrappedf.Message != "failed to read config.yaml" {
		t.Errorf("Message = %v, want %v", wrappedf.Message, "failed to read config.yaml")
	}
}
