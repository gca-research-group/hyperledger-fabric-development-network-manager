package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const migrationGuidance = "use generate --seed <file> followed by validate --output <directory>, or seed generate"

func newRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "experiment-runner",
		Short:         "Generate and validate experiment scenarios",
		Long:          "Generate a complete seed, create a scenario corpus, and validate it in separate calls.\n" + migrationGuidance,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("a subcommand is required: " + migrationGuidance)
		},
	}
	command.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if cmd == command {
			return fmt.Errorf("%w; %s", err, migrationGuidance)
		}
		return err
	})
	command.AddCommand(newGenerateCommand(), newValidateCommand(), newSeedCommand())
	return command
}

func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
