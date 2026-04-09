package cli

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/compose"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
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

		network := compose.ResolveDockerNetworkName(config.Network)

		if err = compose.RemoveContainersInNetwork(network); err != nil {
			return err
		}

		fmt.Println("Containers removed successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(networkDownCmd)

	networkCmd.AddCommand(networkDownCmd)
}
