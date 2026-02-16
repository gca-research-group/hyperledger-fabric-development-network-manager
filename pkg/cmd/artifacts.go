package cmd

import (
	"github.com/spf13/cobra"
)

var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Manage artifacts",
}

func init() {
	rootCmd.AddCommand(artifactsCmd)
}
