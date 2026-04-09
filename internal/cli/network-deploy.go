package cli

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/executor"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/network"
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

		instance := network.NewNetwork(*config, &executor.DefaultExecutor{})

		if err := instance.Deploy(); err != nil {
			return fmt.Errorf("Network deployment failed: %v", err)
		}

		fmt.Println("Deployed successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(networkDeployCmd)

	networkCmd.AddCommand(networkDeployCmd)
}
