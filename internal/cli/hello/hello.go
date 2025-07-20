package hello

import (
	"encoding/json"
	"fmt"

	"github.com/go-cli-template/hello-world-cli/internal/greeting"
	"github.com/spf13/cobra"
)

// Options holds command options
type Options struct {
	IncludeEmoji bool
	JSONOutput   bool
}

// NewCommand creates the hello command
func NewCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "hello",
		Short: "Print a hello world message",
		Long: `Print a hello world message with optional emoji and JSON output.

This command demonstrates a simple greeting without personalization.`,
		Example: `  # Basic hello
  hello-world-cli hello
  
  # With emoji
  hello-world-cli hello --emoji
  
  # JSON output
  hello-world-cli hello --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHello(cmd, opts)
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&opts.IncludeEmoji, "emoji", false, "Include emoji in greeting")
	cmd.Flags().BoolVar(&opts.JSONOutput, "json", false, "Output in JSON format")

	return cmd
}

func runHello(cmd *cobra.Command, opts *Options) error {
	// Create greeting options
	greetOpts := greeting.Options{
		Language:     "en",
		IncludeEmoji: opts.IncludeEmoji,
	}

	// Generate greeting
	greet := greeting.Generate(greetOpts)

	// Output based on format
	if opts.JSONOutput {
		output, err := json.MarshalIndent(greet, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}
		cmd.Println(string(output))
	} else {
		cmd.Println(greet.Message)
	}

	return nil
}
