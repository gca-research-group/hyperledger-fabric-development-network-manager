package cli

import (
	"fmt"
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/directory"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/compose"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/configtx"
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

		if value, _ := directory.IsDirEmpty((*config).Output); value == false {
			return fmt.Errorf("The directory is not empty: %s\n", (*config).Output)
		}

		configTxRenderer := configtx.NewRenderer(config)
		dockerRenderer := compose.NewRenderer(config)

		if err := configTxRenderer.Render(); err != nil {
			return err
		}

		if err := dockerRenderer.Render(); err != nil {
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
