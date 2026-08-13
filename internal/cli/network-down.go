package cli

import (
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
	"github.com/spf13/cobra"
)

var networkDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove all running containers",
	Long:  `Stop and remove all running containers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		if err = workflows.StopNetwork(config); err != nil {
			return err
		}

		slog.Info("Containers removed successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(networkDownCmd)

	networkCmd.AddCommand(networkDownCmd)
}
