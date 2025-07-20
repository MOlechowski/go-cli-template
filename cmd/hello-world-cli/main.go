package main

import (
	"github.com/go-cli-template/hello-world-cli/internal/cli"
	"github.com/go-cli-template/hello-world-cli/internal/errors"
)

// TODO: Replace "hello-world-cli" with your application name
// TODO: Replace "go-cli-template" with your GitHub username/organization

func main() {
	// Set up panic recovery at the top level
	defer errors.PanicHandler()

	// Execute the root command
	if err := cli.Execute(); err != nil {
		// Set debug mode if enabled
		errors.SetDebug(cli.IsDebug())
		// Use the error handler for consistent error presentation
		errors.Exit(err)
	}
}
