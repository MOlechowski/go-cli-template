package errors

import "errors"

// ErrorCode represents a machine-readable error code
type ErrorCode string

// Application error codes
const (
	// General errors
	CodeUnknown        ErrorCode = "UNKNOWN"
	CodeInternal       ErrorCode = "INTERNAL"
	CodeNotImplemented ErrorCode = "NOT_IMPLEMENTED"
	CodeTimeout        ErrorCode = "TIMEOUT"
	CodeCanceled       ErrorCode = "CANCELED"

	// Input/Validation errors
	CodeInvalidInput    ErrorCode = "INVALID_INPUT"
	CodeMissingArgument ErrorCode = "MISSING_ARGUMENT"
	CodeInvalidArgument ErrorCode = "INVALID_ARGUMENT"
	CodeValidation      ErrorCode = "VALIDATION"

	// Configuration errors
	CodeConfig         ErrorCode = "CONFIG"
	CodeConfigNotFound ErrorCode = "CONFIG_NOT_FOUND"
	CodeConfigInvalid  ErrorCode = "CONFIG_INVALID"
	CodeConfigParse    ErrorCode = "CONFIG_PARSE"

	// File/IO errors
	CodeFile           ErrorCode = "FILE"
	CodeFileNotFound   ErrorCode = "FILE_NOT_FOUND"
	CodeFilePermission ErrorCode = "FILE_PERMISSION"
	CodeFileRead       ErrorCode = "FILE_READ"
	CodeFileWrite      ErrorCode = "FILE_WRITE"
	CodeFileCreate     ErrorCode = "FILE_CREATE"

	// Network errors
	CodeNetwork        ErrorCode = "NETWORK"
	CodeNetworkTimeout ErrorCode = "NETWORK_TIMEOUT"
	CodeNetworkDNS     ErrorCode = "NETWORK_DNS"
	CodeNetworkConnect ErrorCode = "NETWORK_CONNECT"

	// Authentication/Authorization
	CodeAuth         ErrorCode = "AUTH"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"

	// Resource errors
	CodeNotFound          ErrorCode = "NOT_FOUND"
	CodeAlreadyExists     ErrorCode = "ALREADY_EXISTS"
	CodeResourceExhausted ErrorCode = "RESOURCE_EXHAUSTED"
)

// ExitCode represents process exit codes following BSD conventions
type ExitCode int

const (
	// ExitSuccess indicates successful completion
	ExitSuccess ExitCode = 0

	// ExitGeneralError indicates a general error
	ExitGeneralError ExitCode = 1

	// ExitMisuse indicates command line usage error
	ExitMisuse ExitCode = 2

	// ExitDataError indicates data format error (BSD code 65)
	ExitDataError ExitCode = 65 // Data format error
	// ExitNoInput indicates cannot open input (BSD code 66)
	ExitNoInput ExitCode = 66 // Cannot open input
	// ExitNoUser indicates user unknown (BSD code 67)
	ExitNoUser ExitCode = 67 // User unknown
	// ExitNoHost indicates host name unknown (BSD code 68)
	ExitNoHost ExitCode = 68 // Host name unknown
	// ExitUnavailable indicates service unavailable (BSD code 69)
	ExitUnavailable ExitCode = 69 // Service unavailable
	// ExitSoftware indicates internal software error (BSD code 70)
	ExitSoftware ExitCode = 70 // Internal software error
	// ExitOSError indicates system error (BSD code 71)
	ExitOSError ExitCode = 71 // System error (e.g., can't fork)
	// ExitOSFile indicates critical OS file missing (BSD code 72)
	ExitOSFile ExitCode = 72 // Critical OS file missing
	// ExitCantCreate indicates can't create output file (BSD code 73)
	ExitCantCreate ExitCode = 73 // Can't create (user) output file
	// ExitIOError indicates input/output error (BSD code 74)
	ExitIOError ExitCode = 74 // Input/output error
	// ExitTempFail indicates temporary failure (BSD code 75)
	ExitTempFail ExitCode = 75 // Temp failure; user is invited to retry
	// ExitProtocol indicates remote protocol error (BSD code 76)
	ExitProtocol ExitCode = 76 // Remote error in protocol
	// ExitNoPerm indicates permission denied (BSD code 77)
	ExitNoPerm ExitCode = 77 // Permission denied
	// ExitConfig indicates configuration error (BSD code 78)
	ExitConfig ExitCode = 78 // Configuration error

	// ExitTimeout indicates command timed out (custom code 124)
	ExitTimeout ExitCode = 124 // Command timed out
	// ExitCanceled indicates command was canceled (custom code 125)
	ExitCanceled ExitCode = 125 // Command was canceled
)

// codeToExit maps error codes to exit codes
var codeToExit = map[ErrorCode]ExitCode{
	// General
	CodeUnknown:        ExitGeneralError,
	CodeInternal:       ExitSoftware,
	CodeNotImplemented: ExitSoftware,
	CodeTimeout:        ExitTimeout,
	CodeCanceled:       ExitCanceled,

	// Input/Validation
	CodeInvalidInput:    ExitMisuse,
	CodeMissingArgument: ExitMisuse,
	CodeInvalidArgument: ExitMisuse,
	CodeValidation:      ExitDataError,

	// Configuration
	CodeConfig:         ExitConfig,
	CodeConfigNotFound: ExitConfig,
	CodeConfigInvalid:  ExitConfig,
	CodeConfigParse:    ExitConfig,

	// File/IO
	CodeFile:           ExitIOError,
	CodeFileNotFound:   ExitNoInput,
	CodeFilePermission: ExitNoPerm,
	CodeFileRead:       ExitIOError,
	CodeFileWrite:      ExitCantCreate,
	CodeFileCreate:     ExitCantCreate,

	// Network
	CodeNetwork:        ExitUnavailable,
	CodeNetworkTimeout: ExitTempFail,
	CodeNetworkDNS:     ExitNoHost,
	CodeNetworkConnect: ExitUnavailable,

	// Auth
	CodeAuth:         ExitNoPerm,
	CodeUnauthorized: ExitNoPerm,
	CodeForbidden:    ExitNoPerm,

	// Resources
	CodeNotFound:          ExitNoInput,
	CodeAlreadyExists:     ExitCantCreate,
	CodeResourceExhausted: ExitUnavailable,
}

// GetExitCode returns the appropriate exit code for an error
func GetExitCode(err error) ExitCode {
	if err == nil {
		return ExitSuccess
	}

	// Check if it's our custom error with a code
	var appErr *Error
	if errors.As(err, &appErr) {
		if exit, ok := codeToExit[appErr.Code]; ok {
			return exit
		}
	}

	// Check specific error types
	if IsValidation(err) {
		return ExitDataError
	}
	if IsConfig(err) {
		return ExitConfig
	}
	if IsFile(err) {
		return ExitIOError
	}
	if IsNetwork(err) {
		return ExitUnavailable
	}

	// Default to general error
	return ExitGeneralError
}

// String returns a human-readable description of the exit code
func (e ExitCode) String() string {
	return getExitCodeString(e)
}

// getExitCodeString returns the string representation for an exit code
// This helper function reduces cyclomatic complexity
func getExitCodeString(e ExitCode) string {
	descriptions := map[ExitCode]string{
		ExitSuccess:      "Success",
		ExitGeneralError: "General error",
		ExitMisuse:       "Command line usage error",
		ExitDataError:    "Data format error",
		ExitNoInput:      "Cannot open input",
		ExitNoUser:       "User unknown",
		ExitNoHost:       "Host name unknown",
		ExitUnavailable:  "Service unavailable",
		ExitSoftware:     "Internal software error",
		ExitOSError:      "System error",
		ExitOSFile:       "Critical OS file missing",
		ExitCantCreate:   "Cannot create output",
		ExitIOError:      "I/O error",
		ExitTempFail:     "Temporary failure",
		ExitProtocol:     "Protocol error",
		ExitNoPerm:       "Permission denied",
		ExitConfig:       "Configuration error",
		ExitTimeout:      "Timeout",
		ExitCanceled:     "Canceled",
	}

	if desc, ok := descriptions[e]; ok {
		return desc
	}
	return "Unknown error"
}
