package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	exampleFlag string
	exampleBool bool
)

var exampleCmd = &cobra.Command{
	Use:   "example [flags]",
	Short: "An example command",
	Long: `This is an example command that demonstrates how to structure
commands in your CLI application.

It shows how to:
- Define flags (string, bool, etc.)
- Handle arguments
- Implement the command logic`,
	Example: `  # Run with default values
  {{.ProjectName}} example

  # Run with custom flag
  {{.ProjectName}} example --flag="custom value"

  # Run with boolean flag
  {{.ProjectName}} example --enable`,
	Args: cobra.MaximumNArgs(1), // Accept at most 1 argument
	RunE: func(cmd *cobra.Command, args []string) error {
		// Command implementation goes here
		fmt.Println("Running example command...")
		
		if exampleFlag != "" {
			fmt.Printf("Flag value: %s\n", exampleFlag)
		}
		
		if exampleBool {
			fmt.Println("Boolean flag is enabled")
		}
		
		if len(args) > 0 {
			fmt.Printf("Argument provided: %s\n", args[0])
		}
		
		// Add your command logic here
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)

	// Local flags - only for this command
	exampleCmd.Flags().StringVarP(&exampleFlag, "flag", "f", "", "An example string flag")
	exampleCmd.Flags().BoolVarP(&exampleBool, "enable", "e", false, "An example boolean flag")
	
	// Mark flag as required (optional)
	// exampleCmd.MarkFlagRequired("flag")
}