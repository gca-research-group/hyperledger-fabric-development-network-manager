package cli

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/executor"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/network"
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

		instance := network.NewContainerManager(*config, &executor.DefaultExecutor{})

		if err := instance.Start(); err != nil {
			return fmt.Errorf("Network starting failed: %v", err)
		}

		fmt.Println("Started successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(networkUpCmd)

	networkCmd.AddCommand(networkUpCmd)
}
