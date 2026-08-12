package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/experiment-runner/generator"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

func TestRunChecksEveryScenario(t *testing.T) {
	outputDirectory := t.TempDir()
	configDirectory := filepath.Join(outputDirectory, "config")
	if err := os.Mkdir(configDirectory, 0755); err != nil {
		t.Fatal(err)
	}

	scenarios := []generator.ScenarioRules{
		{Scenario: "000001", Rules: []validate.RuleID{validate.RuleOrganizationsRequired}},
		{Scenario: "000002", Rules: []validate.RuleID{validate.RuleOrganizationNameRequired}},
	}
	writeScenario(t, configDirectory, "000001", "output: output/example\norganizations: []\n")
	writeScenario(t, configDirectory, "000002", "output: output/example\norganizations:\n  - domain: example.com\n")

	summary, err := Run(outputDirectory, scenarios)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total != 2 || summary.Passed != 2 || summary.Failed != 0 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if _, err := os.Stat(filepath.Join(outputDirectory, "results.json")); err != nil {
		t.Fatalf("results file was not written: %v", err)
	}
}

func writeScenario(t *testing.T, directory, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name+".yaml"), []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}
