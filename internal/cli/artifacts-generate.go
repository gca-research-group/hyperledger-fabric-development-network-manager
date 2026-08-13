package cli

import (
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
	"github.com/spf13/cobra"
)

var artifactsGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate artifacts from a configuration file",
	Long:  `Generate artifacts based on the provided configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var config *config.Config
		var err error

		if config, err = LoadConfig(); err != nil {
			return err
		}

		if err := workflows.GenerateArtifacts(config); err != nil {
			return err
		}

		slog.Info("Artifacts generated successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(artifactsGenerateCmd)

	artifactsCmd.AddCommand(artifactsGenerateCmd)
}
