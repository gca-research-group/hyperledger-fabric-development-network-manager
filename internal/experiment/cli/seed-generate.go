package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/experiment/seed"
	"github.com/spf13/cobra"
)

func newSeedGenerateCommand() *cobra.Command {
	var path string
	var force bool
	command := &cobra.Command{
		Use:   "generate",
		Short: "Write a complete example seed YAML file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeSeed(path, force, cmd.OutOrStdout())
		},
	}
	command.Flags().StringVar(&path, "output", "seed.yaml", "Destination YAML file")
	command.Flags().BoolVar(&force, "force", false, "Overwrite an existing seed file")
	return command
}

func writeSeed(path string, force bool, stdout io.Writer) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create seed directory: %w", err)
	}
	mode := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	if force {
		mode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	file, err := os.OpenFile(path, mode, 0644)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("seed file %q already exists; use --force to overwrite", path)
	}
	if err != nil {
		return fmt.Errorf("create seed file: %w", err)
	}
	_, writeErr := io.WriteString(file, seed.YAML)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("write seed file: %w", err)
	}
	fmt.Fprintf(stdout, "generated seed %s\n", path)
	return nil
}
