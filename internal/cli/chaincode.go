package cli

import (
	"github.com/spf13/cobra"
)

var chaincodeCmd = &cobra.Command{
	Use:   "chaincode",
	Short: "Manage chaincode",
}

func init() {
	rootCmd.AddCommand(chaincodeCmd)
}
