package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/experiment/generator"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/experiment/runner"
)

type options struct {
	seedYAMLPath    string
	mutationCount   int
	outputDirectory string
}

func parseOptions(args []string) (options, error) {
	flags := flag.NewFlagSet("experiment-runner", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var options options
	flags.StringVar(&options.seedYAMLPath, "seed", "", "path to the seed YAML configuration")
	flags.IntVar(&options.mutationCount, "mutation-count", generator.DefaultMutationCount, "number of mutations per scenario")
	flags.StringVar(&options.outputDirectory, "output", "output", "directory for generated scenarios and results")

	if err := flags.Parse(args); err != nil {
		return options, err
	}
	if options.seedYAMLPath == "" {
		return options, errors.New("missing required flag: --seed")
	}
	if flags.NArg() > 0 {
		return options, fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}

	return options, nil
}

func run(args []string, stdout io.Writer) error {
	options, err := parseOptions(args)
	if err != nil {
		return err
	}

	seedYAML, err := os.ReadFile(options.seedYAMLPath)
	if err != nil {
		return fmt.Errorf("read seed YAML %q: %w", options.seedYAMLPath, err)
	}

	generation, err := generator.Generate(seedYAML, options.mutationCount, options.outputDirectory)
	if err != nil {
		return err
	}

	summary, err := runner.RunDirectory(options.outputDirectory)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "generated %d scenarios; processed %d: %d passed, %d partial, %d failed\n", generation.Total, summary.Total, summary.Passed, summary.Partial, summary.Failed)
	if summary.Partial > 0 || summary.Failed > 0 {
		return fmt.Errorf("%d scenarios did not produce all expected validation rules", summary.Partial+summary.Failed)
	}

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}
