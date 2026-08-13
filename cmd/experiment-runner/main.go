package main

import (
	"fmt"
	"log"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/experiment/generator"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/experiment/runner"
)

func main() {
	const outputDirectory = "output"

	generation, err := generator.Generate(outputDirectory, generator.DefaultMutationCount)
	if err != nil {
		log.Fatal(err)
	}

	summary, err := runner.RunDirectory(outputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("generated %d scenarios; processed %d: %d passed, %d partial, %d failed\n", generation.Total, summary.Total, summary.Passed, summary.Partial, summary.Failed)
	if summary.Partial > 0 || summary.Failed > 0 {
		log.Fatalf("%d scenarios did not produce all expected validation rules", summary.Partial+summary.Failed)
	}
}
