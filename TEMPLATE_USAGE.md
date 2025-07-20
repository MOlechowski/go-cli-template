# Go CLI Template Usage Guide

This guide explains how to use this GitHub template to create your own Go CLI project.

## Quick Start

### 1. Create Your Repository

#### Option A: Using GitHub UI
1. Click the "Use this template" button on GitHub
2. Choose owner, repository name, and visibility
3. Clone your new repository locally

#### Option B: Using GitHub CLI
```bash
gh repo create mycli --template MOlechowski/go-cli-template --private --clone
cd mycli
```

### 2. Initialize Your Project

Run the initialization script to automatically customize the template:

```bash
./scripts/init-template.sh
```

The script will:
- Prompt for your project details
- Replace all placeholders throughout the codebase
- Rename directories to match your project name
- Reinitialize the Go module
- Create an initial commit
- Remove itself after completion

## Placeholder Reference

The template uses these placeholders that need to be replaced:

| Placeholder | Description | Example | Used In |
|------------|-------------|---------|----------|
| `{{.ProjectName}}` | Your CLI binary name | `mycli` | go.mod, main.go, commands |
| `{{.GitHubUsername}}` | Your GitHub username | `johndoe` | go.mod, imports |
| `{{.ShortDescription}}` | Brief project description | `A powerful CLI tool` | README.md, root.go |
| `{{.LongDescription}}` | Detailed project description | `This CLI helps developers...` | README.md, root.go |
| `{{.PROJECTNAME}}` | Uppercase name for env vars | `MYCLI` | README.md, config examples |

## Manual Customization (Alternative)

If you prefer to customize manually instead of using the init script:

### 1. Replace Placeholders
Search and replace all placeholders in these files:
- `go.mod`
- `cmd/{{.ProjectName}}/main.go` 
- `internal/cli/*.go`
- `internal/version/*.go`
- `README.md`
- `Taskfile.yml`
- `.air.toml`
- `.gitignore`

### 2. Rename Directories
```bash
mv cmd/{{.ProjectName}} cmd/mycli
```

### 3. Reinitialize Go Module
```bash
rm go.mod go.sum
go mod init github.com/yourusername/mycli
go mod tidy
```

### 4. Update Tests
Run tests to ensure everything works:
```bash
task test
```

## Project Structure

```
.
├── cmd/{{.ProjectName}}/    # Application entry point
│   └── main.go              # Main function
├── internal/                # Private application code
│   ├── cli/                 # Cobra command definitions
│   │   ├── root.go         # Root command and configuration
│   │   ├── version.go      # Version command
│   │   └── example.go      # Example command (remove/modify)
│   └── version/            # Version information
│       └── version.go      # Build-time version injection
├── scripts/                 # Utility scripts
│   └── init-template.sh    # Template initialization
├── .github/                # GitHub configuration
│   ├── workflows/          # GitHub Actions
│   │   ├── ci.yml         # CI pipeline
│   │   └── release.yml    # Release automation
│   └── ISSUE_TEMPLATE/    # Issue templates
├── Taskfile.yml            # Task automation
├── .air.toml              # Hot reload configuration
├── .gitignore             # Git ignore patterns
├── .editorconfig          # Editor configuration
├── go.mod                 # Go module definition
├── README.md              # Project documentation
└── LICENSE                # License file

```

## Development Workflow

### Common Tasks

```bash
# Development with hot reload
task dev

# Run tests
task test

# Lint code
task lint

# Build binary
task build

# Build for all platforms
task build-all

# Create release
task release
```

### Adding New Commands

1. Create a new file in `internal/cli/`:
```go
// internal/cli/mycommand.go
package cli

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

2. Add tests in `internal/cli/mycommand_test.go`

3. Build and test:
```bash
task build
./bin/mycli mycommand
```

## Customization Points

### 1. Remove Example Code
The template includes an example command. Remove it:
- Delete `internal/cli/example.go`
- Delete `internal/cli/example_test.go`
- Update README.md to remove example references

### 2. Configuration
- Global config: Modify `internal/cli/root.go`
- Add new config files: Update `.gitignore`
- Environment variables: Use uppercase project name

### 3. Dependencies
Add dependencies as needed:
```bash
go get github.com/some/package
go mod tidy
```

### 4. CI/CD
- GitHub Actions are pre-configured in `.github/workflows/`
- Modify `ci.yml` for your test requirements
- Update `release.yml` for your release process

### 5. Version Management
Version information is injected at build time:
```bash
task build VERSION=1.0.0
```

## Best Practices

### Code Organization
- Keep business logic in `internal/`
- Commands should be thin wrappers
- Use interfaces for testability
- Follow Go standard project layout

### Testing
- Write tests for all commands
- Use table-driven tests
- Mock external dependencies
- Run tests before commits

### Documentation
- Update README.md with your project specifics
- Document all commands and flags
- Include usage examples
- Keep CHANGELOG.md updated

### Security
- Don't commit sensitive data
- Use environment variables for secrets
- Review dependencies regularly
- Enable GitHub security alerts

### Releases
- Use semantic versioning
- Tag releases properly: `git tag v1.0.0`
- GitHub Actions will create releases automatically
- Include binaries for major platforms

## Troubleshooting

### Init Script Issues
- Ensure you have bash installed
- Run with `bash -x scripts/init-template.sh` for debugging
- Check file permissions: `chmod +x scripts/init-template.sh`

### Build Issues
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to resolve dependencies
- Check for placeholder strings still in code

### Test Failures
- Some tests require `github.com/stretchr/testify`
- Run `go mod tidy` to install test dependencies
- Check for hardcoded values that need updating

## Next Steps

1. Review and update the LICENSE file
2. Configure repository settings on GitHub
3. Set up any required secrets for GitHub Actions
4. Start implementing your CLI logic
5. Update README.md with project-specific information

## Support

For issues with the template itself, please visit:
https://github.com/MOlechowski/go-cli-template/issues