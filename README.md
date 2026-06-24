# Hello World CLI

[![CI](https://github.com/MOlechowski/go-cli-template/actions/workflows/ci.yml/badge.svg)](https://github.com/MOlechowski/go-cli-template/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/MOlechowski/go-cli-template)](https://goreportcard.com/report/github.com/MOlechowski/go-cli-template)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-cli-template/hello-world-cli.svg)](https://pkg.go.dev/github.com/go-cli-template/hello-world-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple hello world CLI demonstrating Go + Cobra with enterprise-ready structure.

## Features

- 🏗️ **Clean, simple structure** - Well-organized code that's easy to understand and extend
- 📊 **Structured logging** - Built-in logging with slog (debug, info, warn, error levels)
- 🌍 **Internationalization** - Support for multiple languages (EN, ES, FR, DE, JA, ZH)
- 📝 **Multiple output formats** - Plain text and JSON output
- 🧪 **Comprehensive testing** - Unit tests with good coverage
- 🔧 **Modern tooling** - automation with mise, hot reload, and CI/CD

## Installation

### Using Go

```bash
go install github.com/go-cli-template/hello-world-cli@latest
```

### Build from Source

```bash
git clone https://github.com/go-cli-template/hello-world-cli.git
cd hello-world-cli
mise run build:default
```

## Usage

### Basic Commands

```bash
# Simple hello world
hello-world-cli hello

# Hello with emoji
hello-world-cli hello --emoji

# Personalized greeting
hello-world-cli greet --name Alice

# Greeting in Spanish
hello-world-cli greet --name Carlos --lang es

# JSON output
hello-world-cli hello --json

# List supported languages
hello-world-cli greet --list-languages

# Enable debug logging
hello-world-cli hello --debug

# Use JSON logging format
hello-world-cli greet --name Alice --log-format=json
```

### Examples

```bash
$ hello-world-cli hello
Hello, World!

$ hello-world-cli hello --emoji
👋 Hello, World!

$ hello-world-cli greet --name Alice --lang fr
Bonjour, Alice!

$ hello-world-cli hello --json
{
  "message": "Hello, World!",
  "language": "en",
  "timestamp": "2024-01-20T10:00:00Z"
}
```

## Development

### Prerequisites

- Go 1.21 or later
- [mise](https://mise.jdx.dev/)

### Common Tasks

```bash
# Run with hot reload
mise run dev:default

# Run tests
mise run test:default

# Build binary
mise run build:default

# Run linter
mise run lint:default

# Build for all platforms
mise run release:build:default

# Format, fix, and lint all code
mise run fix:default
```

### Project Structure

```
├── cmd/hello-world-cli/      # Application entry point
├── internal/                 # Private application code
│   ├── cli/                 # CLI commands
│   │   ├── greet/          # Greet command
│   │   ├── hello/          # Hello command
│   │   └── version/        # Version command
│   └── greeting/           # Core greeting logic
├── pkg/version/            # Public version package
├── scripts/                # Utility scripts
└── docs/                   # Documentation
```

## Customizing This Template

This repository serves as a template for building Go CLI applications. To adapt it for your needs:

1. **Update module name** in `go.mod`
2. **Replace "hello-world-cli"** throughout the codebase with your app name
3. **Modify commands** in `internal/cli/`
4. **Add your business logic** in `internal/`
5. **Update tests** accordingly

See [CUSTOMIZE.md](docs/CUSTOMIZE.md) for detailed instructions.

## Architecture

This CLI follows a clean, simple architecture:

- **Commands** (`internal/cli/`) - CLI command handlers using Cobra
- **Core Logic** (`internal/greeting/`) - Business logic as simple, testable functions
- **Logger** (`internal/logger/`) - Flexible logging system with multiple outputs
- **Public API** (`pkg/version/`) - Exported packages for external use

This structure keeps the code organized and easy to understand while avoiding over-engineering.

See [docs/LOGGING.md](docs/LOGGING.md) for detailed logging documentation.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.