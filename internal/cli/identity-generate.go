package cli

import (
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
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

		if err := workflows.GenerateIdentities(config); err != nil {
			return err
		}

		slog.Info("Identities generated successfully.")
		return nil
	},
}

func init() {
	AddConfigCommand(identityGenerateCmd)

	identityCmd.AddCommand(identityGenerateCmd)
}
