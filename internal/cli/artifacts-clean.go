package cli

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/directory"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
	"github.com/spf13/cobra"
)

var artifactsCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the output folder",
	Long:  `Remove all files from the output folder.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var config *config.Config
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
	AddConfigCommand(artifactsCleanCmd)

	artifactsCmd.AddCommand(artifactsCleanCmd)
}
