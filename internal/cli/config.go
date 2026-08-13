package cli

import (
	"errors"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/application"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/config"
	"github.com/spf13/cobra"
)

var configPath string
var workflows = application.NewWorkflows(nil)

func LoadConfig() (*config.Config, error) {
	if configPath == "" {
		return nil, errors.New("missing required flag: --config")
	}

	return config.LoadConfigFromPath(configPath)
}

func AddConfigCommand(command *cobra.Command) {
	command.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		"Path to configuration file",
	)

	command.MarkFlagFilename("config")
	command.MarkFlagRequired("config")
}
