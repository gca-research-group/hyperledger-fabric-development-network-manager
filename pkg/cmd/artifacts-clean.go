package cmd

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/directory"
	"github.com/spf13/cobra"
)

var artifactsCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the output folder",
	Long:  `Remove all files from the output folder.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var config *pkg.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		directory.RemoveFolderIfExists(config.Output)

		fmt.Println("Artifacts removed successfully.")
		return nil
	},
}

func init() {
	artifactsCmd.AddCommand(artifactsCleanCmd)

	artifactsCleanCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	artifactsCleanCmd.MarkFlagRequired("config")
}
