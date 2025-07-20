package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/go-cli-template/hello-world-cli/internal/cli/greet"
	"github.com/go-cli-template/hello-world-cli/internal/cli/hello"
	versioncmd "github.com/go-cli-template/hello-world-cli/internal/cli/version"
	"github.com/go-cli-template/hello-world-cli/internal/domain/greeting"
	"github.com/go-cli-template/hello-world-cli/internal/domain/language"
	"github.com/go-cli-template/hello-world-cli/pkg/version"
)

var (
	cfgFile string
	verbose bool
)

// TODO: Replace "hello-world-cli" with your application name throughout this file

var rootCmd = &cobra.Command{
	Use:   "hello-world-cli",
	Short: "A simple hello world CLI demonstrating Go + Cobra",
	Long: `Hello World CLI is a demonstration of building a well-structured
command-line application in Go using the Cobra framework.

This CLI showcases:
- Enterprise-ready directory structure
- Domain-driven design principles
- Multiple commands with sub-commands
- Internationalization support
- JSON output formatting
- Comprehensive testing approach`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Initialize services (in a real app, this might use dependency injection)
	languageService := language.NewService()
	greetingService := greeting.NewService(languageService)

	// Add commands
	rootCmd.AddCommand(hello.NewCommand(greetingService))
	rootCmd.AddCommand(greet.NewCommand(greetingService))
	rootCmd.AddCommand(versioncmd.NewCommand())

	// Persistent flags - global for all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hello-world-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add version info
	rootCmd.Version = version.String()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hello-world-cli") // TODO: Replace with your app name
	}

	// Set environment variable prefix
	viper.SetEnvPrefix("HELLO_WORLD_CLI") // TODO: Replace with your app name in uppercase
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}