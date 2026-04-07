package cli

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
	"github.com/spf13/cobra"
)

var identityGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate identities",
	Long:  `Generate identities.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		containerManager := network.NewContainerManager(*config, &executor.DefaultExecutor{})

		if err := containerManager.RunCAContainers(); err != nil {
			return err
		}

		instance := network.NewIdentityManager(*config, &executor.DefaultExecutor{})

		if err := instance.GenerateAll(); err != nil {
			return fmt.Errorf("Generation of identities have failed: %v", err)
		}

		if err := containerManager.StopCertificateAuthorities(); err != nil {
			return err
		}

		fmt.Println("Identities generated successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(identityGenerateCmd)

	identityCmd.AddCommand(identityGenerateCmd)
}
