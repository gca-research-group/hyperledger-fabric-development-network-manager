package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fno",
	Short: "fno generates and manages artifacts",
	Long:  "fno is a CLI for generating and managing artifacts from configuration files.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
