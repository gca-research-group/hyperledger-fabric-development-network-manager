package cli

import (
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
	"github.com/spf13/cobra"
)

var networkDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the network",
	Long:  `Deploy the network starting the containers, generating the identities, generating the genesis block, creating the channels, and joining orderers and peers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		if err := workflows.DeployNetwork(config); err != nil {
			return err
		}

		slog.Info("Deployed successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(networkDeployCmd)

	networkCmd.AddCommand(networkDeployCmd)
}
