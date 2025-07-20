package cli

import (
	"fmt"
	"os"

	"github.com/go-cli-template/hello-world-cli/internal/cli/greet"
	"github.com/go-cli-template/hello-world-cli/internal/cli/hello"
	versioncmd "github.com/go-cli-template/hello-world-cli/internal/cli/version"
	"github.com/go-cli-template/hello-world-cli/internal/logger"
	"github.com/go-cli-template/hello-world-cli/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	verbose   bool
	debug     bool
	logLevel  string
	logFormat string
)

// TODO: Replace "hello-world-cli" with your application name throughout this file

var rootCmd = &cobra.Command{
	Use:   "hello-world-cli",
	Short: "A simple hello world CLI demonstrating Go + Cobra",
	Long: `Hello World CLI is a demonstration of building a well-structured
command-line application in Go using the Cobra framework.

This CLI showcases:
- Clean, simple architecture
- Structured logging with slog
- Multiple commands with sub-commands
- Internationalization support
- JSON output formatting
- Comprehensive testing approach`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Configure logger based on viper configuration and flags
		cfg := logger.DefaultConfig()

		// Read from viper configuration first
		if viper.IsSet("log.level") {
			cfg.Level = viper.GetString("log.level")
		}
		if viper.IsSet("log.format") {
			cfg.Format = viper.GetString("log.format")
		}

		// Command line flags override config file
		if debug {
			cfg.Level = "debug"
		} else if logLevel != "" {
			cfg.Level = logLevel
		}

		if logFormat != "" {
			cfg.Format = logFormat
		}

		// Environment variables for logging
		if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
			cfg.Level = envLevel
		}
		if envFormat := os.Getenv("LOG_FORMAT"); envFormat != "" {
			cfg.Format = envFormat
		}

		// Create and set the logger
		log := logger.New(cfg)
		logger.SetDefault(log)

		// Add logger to context
		ctx := logger.WithContext(cmd.Context(), log)
		cmd.SetContext(ctx)

		// Log startup information at debug level
		log.Debug("starting hello-world-cli",
			"version", version.Version,
			"args", os.Args[1:],
			"log_level", cfg.Level,
			"log_format", cfg.Format,
		)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// IsDebug returns whether debug mode is enabled
func IsDebug() bool {
	return debug
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add commands
	rootCmd.AddCommand(hello.NewCommand())
	rootCmd.AddCommand(greet.NewCommand())
	rootCmd.AddCommand(versioncmd.NewCommand())

	// Persistent flags - global for all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hello-world-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging (includes file:line info)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "set log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "set log format (text, json)")

	// Bind flags to viper
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flag: %v\n", err)
	}
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flag: %v\n", err)
	}
	if err := viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flag: %v\n", err)
	}
	if err := viper.BindPFlag("log.format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding flag: %v\n", err)
	}

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
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
