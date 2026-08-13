package cli

import (
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
	"github.com/spf13/cobra"
)

var networkUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the network",
	Long:  `Start the network containers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		if err := workflows.StartNetwork(config); err != nil {
			return err
		}

		slog.Info("Started successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(networkUpCmd)

	networkCmd.AddCommand(networkUpCmd)
}
