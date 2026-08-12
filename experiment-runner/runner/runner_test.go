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
		{Scenario: "000001", Rules: []validate.RuleID{validate.RuleOrganizationsRequired, validate.RuleOrdererTopologyRequired}},
		{Scenario: "000002", Rules: []validate.RuleID{validate.RuleOrganizationsRequired, validate.RuleApplicationCapabilityUnsupported}},
		{Scenario: "000003", Rules: []validate.RuleID{validate.RuleChannelNameInvalid}},
	}
	writeScenario(t, configDirectory, "000001", "output: output/example\norganizations: []\n")
	writeScenario(t, configDirectory, "000002", "output: output/example\ncapabilities:\n  channel: V2_0\n  application: V2_5\n  orderer: V2_0\norganizations: []\n")
	writeScenario(t, configDirectory, "000003", "output: output/example\norganizations: []\n")

	summary, err := Run(outputDirectory, scenarios)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total != 3 || summary.Passed != 1 || summary.Partial != 1 || summary.Failed != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if summary.Results[0].Status != StatusPassed || summary.Results[1].Status != StatusPartial || summary.Results[2].Status != StatusFailed {
		t.Fatalf("unexpected result states: %+v", summary.Results)
	}
	if len(summary.Results[1].Missing) != 1 || summary.Results[1].Missing[0] != validate.RuleApplicationCapabilityUnsupported {
		t.Fatalf("unexpected missing rules: %+v", summary.Results[1].Missing)
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
