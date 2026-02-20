package cli

import (
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage network",
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
