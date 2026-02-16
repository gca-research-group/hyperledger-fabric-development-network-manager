package cmd

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/fabric"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the network",
	Long:  `Deploy the network starting the containers, generating the identities, generating the genesis block, creating the channels, and joining orderers and peers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *pkg.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		instance, err := fabric.NewFabric(*config, &fabric.DefaultExecutor{})

		if err != nil {
			return err
		}

		if err := instance.DeployNetwork(); err != nil {
			return fmt.Errorf("Network deployment failed: %v", err)
		}

		fmt.Println("Deployed successfully.")
		return nil
	},
}

func init() {
	deployCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	deployCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(deployCmd)
}
