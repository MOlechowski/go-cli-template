package errors

import (
	"fmt"
	"testing"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode ExitCode
	}{
		{
			name:     "nil error returns success",
			err:      nil,
			wantCode: ExitSuccess,
		},
		{
			name:     "validation error returns data error",
			err:      &ValidationError{Field: "test", Message: "invalid"},
			wantCode: ExitDataError,
		},
		{
			name:     "config error returns config exit code",
			err:      &ConfigError{Key: "test", Message: "invalid"},
			wantCode: ExitConfig,
		},
		{
			name:     "file error returns IO error",
			err:      &FileError{Path: "test", Operation: "read"},
			wantCode: ExitIOError,
		},
		{
			name:     "network error returns unavailable",
			err:      &NetworkError{URL: "test", Operation: "GET"},
			wantCode: ExitUnavailable,
		},
		{
			name:     "error with timeout code",
			err:      &Error{Code: CodeTimeout},
			wantCode: ExitTimeout,
		},
		{
			name:     "error with permission code",
			err:      &Error{Code: CodeFilePermission},
			wantCode: ExitNoPerm,
		},
		{
			name:     "unknown error returns general error",
			err:      fmt.Errorf("unknown error"),
			wantCode: ExitGeneralError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExitCode(tt.err); got != tt.wantCode {
				t.Errorf("GetExitCode() = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestExitCodeString(t *testing.T) {
	tests := []struct {
		code ExitCode
		want string
	}{
		{ExitSuccess, "Success"},
		{ExitGeneralError, "General error"},
		{ExitMisuse, "Command line usage error"},
		{ExitDataError, "Data format error"},
		{ExitNoInput, "Cannot open input"},
		{ExitNoUser, "User unknown"},
		{ExitNoHost, "Host name unknown"},
		{ExitUnavailable, "Service unavailable"},
		{ExitSoftware, "Internal software error"},
		{ExitOSError, "System error"},
		{ExitOSFile, "Critical OS file missing"},
		{ExitCantCreate, "Cannot create output"},
		{ExitIOError, "I/O error"},
		{ExitTempFail, "Temporary failure"},
		{ExitProtocol, "Protocol error"},
		{ExitNoPerm, "Permission denied"},
		{ExitConfig, "Configuration error"},
		{ExitTimeout, "Timeout"},
		{ExitCanceled, "Canceled"},
		{ExitCode(999), "Unknown error"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.code.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeToExitMapping(t *testing.T) {
	// Test that all error codes have a mapping
	codes := []ErrorCode{
		CodeUnknown, CodeInternal, CodeTimeout, CodeCanceled,
		CodeInvalidInput, CodeMissingArgument, CodeValidation,
		CodeConfig, CodeConfigNotFound, CodeConfigInvalid,
		CodeFile, CodeFileNotFound, CodeFilePermission,
		CodeNetwork, CodeNetworkTimeout, CodeNetworkConnect,
		CodeAuth, CodeUnauthorized, CodeForbidden,
		CodeNotFound, CodeAlreadyExists,
	}

	for _, code := range codes {
		err := &Error{Code: code}
		exitCode := GetExitCode(err)
		if exitCode == ExitGeneralError && code != CodeUnknown {
			t.Errorf("Code %v maps to general error, expected specific exit code", code)
		}
	}
}
