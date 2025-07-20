# {{.ProjectName}}

{{.ShortDescription}}

## Installation

### Using Go

```bash
go install github.com/{{.GitHubUsername}}/{{.ProjectName}}@latest
```

### Download Binary

Download the latest release from the [releases page](https://github.com/{{.GitHubUsername}}/{{.ProjectName}}/releases).

### Build from Source

```bash
git clone https://github.com/{{.GitHubUsername}}/{{.ProjectName}}.git
cd {{.ProjectName}}
task build
```

## Usage

```bash
{{.ProjectName}} [command] [flags]
```

### Commands

- `version` - Display version information
- `example` - Example command (customize or remove)
- `help` - Display help for any command

### Global Flags

- `--config string` - Config file (default is $HOME/.{{.ProjectName}}.yaml)
- `-v, --verbose` - Enable verbose output
- `-h, --help` - Display help

### Examples

```bash
# Display version
{{.ProjectName}} version

# Run example command
{{.ProjectName}} example --flag="value"

# Use custom config
{{.ProjectName}} --config=/path/to/config.yaml example
```

## Development

This project uses [Task](https://taskfile.dev/) for build automation.

### Prerequisites

- Go 1.21 or later
- Task (see [installation instructions](https://taskfile.dev/installation/))

### Common Tasks

```bash
# Run the application
task run -- [args]

# Run with hot reload
task dev

# Run tests
task test

# Run linter
task lint

# Build for current platform
task build

# Build for all platforms
task build-all

# Create release
task release
```

### Project Structure

```
.
├── cmd/
│   └── {{.ProjectName}}/
│       └── main.go          # Application entry point
├── internal/
│   ├── cli/                 # CLI commands
│   │   ├── root.go         # Root command
│   │   ├── version.go      # Version command
│   │   └── example.go      # Example command
│   └── version/            # Version information
│       └── version.go
├── .github/                # GitHub Actions workflows
├── Taskfile.yml           # Task automation
├── go.mod                 # Go module definition
└── README.md             # This file
```

### Adding a New Command

1. Create a new file in `internal/cli/`:

```go
package cli

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  `Detailed description`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    // Add flags
    myCmd.Flags().StringP("flag", "f", "", "Flag description")
}
```

2. The command will be automatically available when you rebuild.

## Configuration

{{.ProjectName}} supports configuration via:

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

### Configuration File

Create a config file at `~/.{{.ProjectName}}.yaml`:

```yaml
# Example configuration
verbose: true
# Add your configuration options here
```

### Environment Variables

All configuration options can be set via environment variables. The pattern is:
- Prefix with your project name in uppercase
- Replace dots with underscores

Example:
```bash
export {{.PROJECTNAME}}_VERBOSE=true
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Standards

- Write tests for new functionality
- Ensure all tests pass (`task test`)
- Run linter and fix issues (`task lint`)
- Update documentation as needed

## License

[Choose a license and add it here]

## Using This Template

This repository is a GitHub template. To use it:

1. Click "Use this template" on GitHub
2. Clone your new repository
3. Run the initialization script:
   ```bash
   ./scripts/init-template.sh
   ```
4. Follow the prompts to customize the template
5. Start building your CLI!

### Template Placeholders

The following placeholders should be replaced:
- `{{.ProjectName}}` - Your project name (e.g., `mycli`)
- `{{.GitHubUsername}}` - Your GitHub username
- `{{.ShortDescription}}` - Brief project description
- `{{.LongDescription}}` - Detailed project description
- `{{.PROJECTNAME}}` - Uppercase project name for env vars