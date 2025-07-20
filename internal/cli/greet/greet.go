package greet

import (
	"encoding/json"
	"fmt"

	"github.com/go-cli-template/hello-world-cli/internal/greeting"
	"github.com/spf13/cobra"
)

// Options holds command options
type Options struct {
	Name         string
	Language     string
	IncludeEmoji bool
	JSONOutput   bool
	ListLangs    bool
}

// NewCommand creates the greet command
func NewCommand() *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "greet",
		Short: "Print a personalized greeting",
		Long: `Print a personalized greeting with support for multiple languages.

This command demonstrates personalized greetings with internationalization support.`,
		Example: `  # Basic greeting
  hello-world-cli greet --name Alice
  
  # Spanish greeting with emoji
  hello-world-cli greet --name Carlos --lang es --emoji
  
  # List supported languages
  hello-world-cli greet --list-languages`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGreet(cmd, opts)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name to greet")
	cmd.Flags().StringVarP(&opts.Language, "lang", "l", "en", "Language code (en, es, fr, de, ja, zh)")
	cmd.Flags().BoolVar(&opts.IncludeEmoji, "emoji", false, "Include emoji in greeting")
	cmd.Flags().BoolVar(&opts.JSONOutput, "json", false, "Output in JSON format")
	cmd.Flags().BoolVar(&opts.ListLangs, "list-languages", false, "List all supported languages")

	return cmd
}

func runGreet(cmd *cobra.Command, opts *Options) error {
	// Handle list languages request
	if opts.ListLangs {
		langs := greeting.GetSupportedLanguages()
		cmd.Println("Supported languages:")
		for _, lang := range langs {
			cmd.Printf("  %s\n", lang)
		}
		return nil
	}

	// Validate name is provided
	if opts.Name == "" {
		return fmt.Errorf("name is required (use --name flag)")
	}

	// Create greeting options
	greetOpts := greeting.Options{
		Name:         opts.Name,
		Language:     opts.Language,
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
