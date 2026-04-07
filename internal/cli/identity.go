package cli

import (
	"github.com/spf13/cobra"
)

var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Generate identities",
}

func init() {
	rootCmd.AddCommand(identityCmd)
}
