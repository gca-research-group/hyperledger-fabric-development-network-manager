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

	fmt.Printf("generated %d scenarios; %d passed, %d failed\n", summary.Total, summary.Passed, summary.Failed)
	if summary.Failed > 0 {
		log.Fatalf("%d scenarios did not produce an expected validation rule", summary.Failed)
	}
}
