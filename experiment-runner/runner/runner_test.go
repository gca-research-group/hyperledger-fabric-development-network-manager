package runner

import (
	"encoding/json"
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
	manifest, err := json.Marshal(scenarios)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDirectory, "scenarios.json"), manifest, 0644); err != nil {
		t.Fatal(err)
	}

	summary, err := RunDirectory(outputDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total != 3 || summary.Passed != 1 || summary.Partial != 1 || summary.Failed != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	data, err := os.ReadFile(filepath.Join(outputDirectory, "results.json"))
	if err != nil {
		t.Fatalf("results file was not written: %v", err)
	}
	var document struct {
		Summary
		Results []Result `json:"results"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatalf("decode results: %v", err)
	}
	if document.Results[0].Status != StatusPassed || document.Results[1].Status != StatusPartial || document.Results[2].Status != StatusFailed {
		t.Fatalf("unexpected result states: %+v", document.Results)
	}
	if len(document.Results[1].Missing) != 1 || document.Results[1].Missing[0] != validate.RuleApplicationCapabilityUnsupported {
		t.Fatalf("unexpected missing rules: %+v", document.Results[1].Missing)
	}
}

func writeScenario(t *testing.T, directory, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, name+".yaml"), []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}
