package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/{{.GitHubUsername}}/{{.ProjectName}}/internal/version"
)

var (
	jsonOutput bool
	shortVer   bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print detailed version information about {{.ProjectName}}.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		info := version.GetBuildInfo()

		if shortVer {
			fmt.Println(info.Version)
			return nil
		}

		if jsonOutput {
			data, err := json.MarshalIndent(info, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal version info: %w", err)
			}
			fmt.Println(string(data))
		} else {
			fmt.Printf("{{.ProjectName}} version %s\n", info.Version)
			fmt.Printf("  Build Time: %s\n", info.BuildTime)
			if info.GitCommit != "" {
				fmt.Printf("  Git Commit: %s\n", info.GitCommit)
			}
			fmt.Printf("  Go Version: %s\n", info.GoVersion)
			fmt.Printf("  Platform:   %s\n", info.Platform)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version in JSON format")
	versionCmd.Flags().BoolVar(&shortVer, "short", false, "Print just the version number")
}