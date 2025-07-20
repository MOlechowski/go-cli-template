# CLAUDE.md - Project-Specific Configuration for go-cli-template

This file provides project-specific guidance and configuration for Claude Code when working with this Go CLI template.

## Project Overview

This is a Go CLI template project that provides a foundation for building command-line applications with:
- Structured logging with slog
- Comprehensive error handling
- Command structure using Cobra
- Configuration management with Viper
- Testing patterns and examples

## Linked Rules and Guidelines

@RULES-GO.md - Go-specific coding standards and best practices for this project

## Build & Development Commands

```bash
# Build the project
go build -o hello-world-cli cmd/hello-world-cli/main.go

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run linter
golangci-lint run

# Format code
gofmt -w .
goimports -w .

# Run the CLI
go run cmd/hello-world-cli/main.go --help
```

## Project Structure

```
go-cli-template/
├── cmd/                    # Main applications
│   └── hello-world-cli/    # CLI entry point
├── internal/              # Private application code
│   ├── cli/              # Command implementations
│   ├── errors/           # Error handling system
│   ├── greeting/         # Business logic
│   └── logger/           # Logging system
├── pkg/                   # Public packages
│   └── version/          # Version information
├── docs/                  # Documentation
├── .golangci.yml         # Linter configuration
└── RULES-GO.md           # Go coding standards

## Key Patterns in This Project

### Error Handling
- Uses custom error types in `internal/errors`
- Provides user-friendly error messages
- Implements proper exit codes for scripting

### Logging
- Structured logging with slog
- Automatic file:line info in debug mode
- Configurable via environment, flags, or config file

### Command Structure
- Commands organized under `internal/cli`
- Each command in its own package
- Shared configuration through root command

## Testing Guidelines

- Table-driven tests are preferred
- Test files next to implementation
- Use `testdata/` for test fixtures
- Mock external dependencies

## Common Tasks

### Adding a New Command
1. Create new package under `internal/cli/`
2. Implement command with `cobra.Command`
3. Register in `internal/cli/root.go`
4. Add tests

### Modifying Error Handling
- Error types defined in `internal/errors/errors.go`
- Exit codes in `internal/errors/codes.go`
- Error presentation in `internal/errors/handler.go`

### Updating Logging
- Logger configuration in `internal/logger/logger.go`
- Color handler for terminal output
- Context-based logger propagation

## Important Notes

- Follow RULES-GO.md for all Go code in this project
- Maintain backward compatibility in pkg/ directory
- Keep internal/ packages private
- Update tests when modifying functionality
- Run linter before committing