package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func newSeedCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "seed",
		Short: "Manage experiment seeds",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("seed requires a subcommand: experiment-runner seed generate [--output seed.yaml] [--force]")
		},
	}
	command.AddCommand(newSeedGenerateCommand())
	return command
}
