package version

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/go-cli-template/hello-world-cli/pkg/version"
)

// Options holds command options
type Options struct {
	JSONOutput bool
	Short      bool
}

// NewCommand creates the version command
func NewCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print detailed version information about hello-world-cli.`,
		Example: `  # Show version info
  hello-world-cli version
  
  # Show short version
  hello-world-cli version --short
  
  # Show version in JSON format
  hello-world-cli version --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(opts)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&opts.JSONOutput, "json", false, "Output version in JSON format")
	cmd.Flags().BoolVar(&opts.Short, "short", false, "Print just the version number")

	return cmd
}

func runVersion(opts *Options) error {
	info := version.GetBuildInfo()

	if opts.Short {
		fmt.Println(info.Version)
		return nil
	}

	if opts.JSONOutput {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Printf("hello-world-cli version %s\n", info.Version)
		fmt.Printf("  Build Time: %s\n", info.BuildTime)
		if info.GitCommit != "" {
			fmt.Printf("  Git Commit: %s\n", info.GitCommit)
		}
		fmt.Printf("  Go Version: %s\n", info.GoVersion)
		fmt.Printf("  Platform:   %s\n", info.Platform)
	}

	return nil
}