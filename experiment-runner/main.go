package main

import (
	"fmt"
	"log"

	"github.com/gca-research-group/fabric-network-orchestrator/experiment-runner/generator"
	"github.com/gca-research-group/fabric-network-orchestrator/experiment-runner/runner"
)

func main() {
	const outputDirectory = "output"

	_, err := generator.Generate(outputDirectory, generator.DefaultMutationCount)
	if err != nil {
		log.Fatal(err)
	}

	summary, err := runner.RunDirectory(outputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("generated %d scenarios; %d passed, %d partial, %d failed\n", summary.Total, summary.Passed, summary.Partial, summary.Failed)
	if summary.Partial > 0 || summary.Failed > 0 {
		log.Fatalf("%d scenarios did not produce all expected validation rules", summary.Partial+summary.Failed)
	}
}
