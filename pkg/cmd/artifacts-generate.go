package cmd

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/docker"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/directory"
	"github.com/spf13/cobra"
)

var force bool

var artifactsGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate artifacts from a configuration file",
	Long:  `Generate artifacts based on the provided configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var config *pkg.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		if force {
			directory.RemoveFolderIfExists(config.Output)
		}

		if value, _ := directory.IsDirEmpty((*config).Output); value == false {
			return fmt.Errorf("The directory is not empty: %s\n", (*config).Output)
		}

		configTxRenderer := configtx.NewRenderer(config)
		dockerRenderer := docker.NewRenderer(config)

		if err := configTxRenderer.Render(); err != nil {
			return err
		}

		if err := dockerRenderer.Render(); err != nil {
			return err
		}

		fmt.Println("Artifacts generated successfully.")
		return nil
	},
}

func init() {
	artifactsCmd.AddCommand(artifactsGenerateCmd)

	artifactsGenerateCmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	artifactsGenerateCmd.MarkFlagRequired("config")

	artifactsGenerateCmd.Flags().BoolVarP(
		&force,
		"force",
		"f",
		false,
		"Remove existing files from the output folder",
	)

}
