package cmd

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/docker"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/constants"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove all running containers",
	Long:  `Stop and remove all running containers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *pkg.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		network := config.Network

		if network == "" {
			network = constants.DEFAULT_NETORK
		}

		if err = docker.RemoveContainersInNetwork(network); err != nil {
			return err
		}

		fmt.Println("Containers removed successfully.")
		return nil
	},
}

func init() {
	downCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	downCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(downCmd)
}
