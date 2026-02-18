package cli

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/fabric"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the network",
	Long:  `Start the network containers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		instance, err := fabric.NewFabric(*config, &fabric.DefaultExecutor{})

		if err != nil {
			return err
		}

		if err := instance.Start(); err != nil {
			return fmt.Errorf("Network starting failed: %v", err)
		}

		fmt.Println("Started successfully.")
		return nil
	},
}

func init() {
	upCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	upCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(upCmd)
}
