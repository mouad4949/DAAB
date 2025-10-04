package main

import (
	"fmt"
	"os"

	initcmd "github.com/mouad4949/DAAB/internal/init"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "daab",
		Short: "DAAB - Your DevOps Automation Buddy",
		Long: `DAAB is a CLI tool that automates deployment workflows.
It helps you deploy your applications to the cloud with zero friction.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	rootCmd.AddCommand(initcmd.NewInitCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
