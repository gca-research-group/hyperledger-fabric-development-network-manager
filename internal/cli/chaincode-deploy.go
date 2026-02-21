package cli

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/chaincode"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
	"github.com/spf13/cobra"
)

var chaincodeDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Stop and remove all running containers",
	Long:  `Stop and remove all running containers.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		c := chaincode.NewChaincode(config, &executor.DefaultExecutor{})

		if err := c.Publish(); err != nil {
			return fmt.Errorf("Chaincode deployment failed: %v", err)
		}

		fmt.Println("Chaincodes deploied successfully.")
		return nil
	},
}

func init() {
	chaincodeDeployCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	chaincodeDeployCmd.MarkFlagRequired("config")

	chaincodeCmd.AddCommand(chaincodeDeployCmd)
}
