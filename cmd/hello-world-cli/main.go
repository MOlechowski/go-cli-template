package main

import (
	"fmt"
	"os"

	"github.com/go-cli-template/hello-world-cli/internal/cli"
)

// TODO: Replace "hello-world-cli" with your application name
// TODO: Replace "go-cli-template" with your GitHub username/organization

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}