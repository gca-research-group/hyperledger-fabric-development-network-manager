package cli

import (
	"fmt"
	"os"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fno",
	Short: "fno generates and manages artifacts",
	Long:  "fno is a CLI for generating and managing artifacts from configuration files.",
}

func Execute() {
	if err := logger.Setup(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
