package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tool",
	Short: "Tool generates and manages artifacts",
	Long:  "Tool is a CLI for generating and managing artifacts from configuration files.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
