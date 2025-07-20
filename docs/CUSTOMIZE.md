# Customization Guide

This guide explains how to adapt the Hello World CLI template for your own project.

## Quick Start

### 1. Use as GitHub Template

1. Click "Use this template" on GitHub
2. Clone your new repository
3. Follow the customization steps below

### 2. Search and Replace

Replace these strings throughout the codebase:

| Find | Replace With | Description |
|------|--------------|-------------|
| `hello-world-cli` | `your-app-name` | Binary and project name |
| `go-cli-template` | `your-username` | GitHub username/org |
| `Hello World CLI` | `Your App Name` | Display name |
| `Hello, World!` | Your default message | Default output |

**TODO markers**: Search for `TODO:` comments in the code for specific customization points.

### 3. Update Key Files

1. **go.mod**: Update module path
   ```go
   module github.com/your-username/your-app-name
   ```

2. **Taskfile.yml**: Update binary name and paths
3. **.gitignore**: Update app-specific config file names
4. **README.md**: Update with your project information

## Adding Your Own Commands

### 1. Create Command File

Create a new file in `internal/cli/yourcommand/`:

```go
package yourcommand

import (
    "github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "yourcommand",
        Short: "Brief description",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
            return nil
        },
    }
    
    // Add flags
    cmd.Flags().StringP("flag", "f", "", "Flag description")
    
    return cmd
}
```

### 2. Register Command

In `internal/cli/root.go`, add:

```go
import "github.com/your-username/your-app/internal/cli/yourcommand"

func init() {
    // ... existing code ...
    rootCmd.AddCommand(yourcommand.NewCommand())
}
```

## Adding Business Logic

### 1. Create Domain Service

Create service in `internal/domain/yourdomain/`:

```go
package yourdomain

type Service struct {
    // dependencies
}

func NewService() *Service {
    return &Service{}
}

func (s *Service) YourMethod() error {
    // Business logic here
    return nil
}
```

### 2. Wire Up Dependencies

In `internal/cli/root.go`:

```go
yourService := yourdomain.NewService()
rootCmd.AddCommand(yourcommand.NewCommand(yourService))
```

## Directory Structure

- **cmd/**: Keep minimal, just entry point
- **internal/cli/**: CLI commands only, no business logic
- **internal/domain/**: All business logic goes here
- **internal/infrastructure/**: External integrations (DB, APIs)
- **pkg/**: Public packages that others can import

## Testing

Add tests alongside your code:
- `yourcommand.go` → `yourcommand_test.go`
- `service.go` → `service_test.go`

Run tests with: `task test`

## CI/CD

The template includes GitHub Actions that will:
- Run tests on multiple OS
- Run linter
- Build binaries
- Create releases on tags

Just push tags like `v1.0.0` to create releases.

## Common Customizations

### Add Database

1. Add driver to `go.mod`
2. Create `internal/infrastructure/database/`
3. Add connection logic
4. Inject into services

### Add Configuration

1. Create config struct in `internal/infrastructure/config/`
2. Use Viper to load from file/env
3. Pass config to services

### Add API Client

1. Create client in `internal/infrastructure/api/`
2. Add retry/timeout logic
3. Inject into services

## Removing Example Code

To clean up the hello world example:

1. Delete `internal/cli/hello/` and `internal/cli/greet/`
2. Delete `internal/domain/greeting/` and `internal/domain/language/`
3. Remove command registrations from `internal/cli/root.go`
4. Update tests

## Tips

- Start small, add complexity as needed
- Keep business logic in domain layer
- Use dependency injection
- Write tests as you go
- Follow Go conventions