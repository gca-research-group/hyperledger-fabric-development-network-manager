package cli

import (
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
	"github.com/spf13/cobra"
)

var imagePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull the docker images",
	Long:  `Pull the docker images.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		if err = workflows.PullImages(config); err != nil {
			return err
		}

		slog.Info("Images pulled successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(imagePullCmd)

	imageCmd.AddCommand(imagePullCmd)
}
