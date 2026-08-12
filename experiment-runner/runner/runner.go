package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gca-research-group/fabric-network-orchestrator/experiment-runner/generator"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

type Result struct {
	Scenario string            `json:"scenario"`
	Expected []validate.RuleID `json:"expectedRules"`
	Actual   []validate.RuleID `json:"actualRules,omitempty"`
	Passed   bool              `json:"passed"`
	Error    string            `json:"error,omitempty"`
}

type Summary struct {
	Total   int      `json:"total"`
	Passed  int      `json:"passed"`
	Failed  int      `json:"failed"`
	Results []Result `json:"results"`
}

// RunDirectory reads the generated manifest and runs its complete corpus.
func RunDirectory(outputDirectory string) (Summary, error) {
	data, err := os.ReadFile(filepath.Join(outputDirectory, "scenarios.json"))
	if err != nil {
		return Summary{}, fmt.Errorf("read scenario manifest: %w", err)
	}

	var scenarios []generator.ScenarioRules
	if err := json.Unmarshal(data, &scenarios); err != nil {
		return Summary{}, fmt.Errorf("decode scenario manifest: %w", err)
	}

	return Run(outputDirectory, scenarios)
}

// Run validates every generated configuration and checks that the reported rule
// is one of the mutations recorded for that scenario.
func Run(outputDirectory string, scenarios []generator.ScenarioRules) (Summary, error) {
	summary := Summary{Total: len(scenarios), Results: make([]Result, 0, len(scenarios))}

	for _, scenario := range scenarios {
		result := Result{Scenario: scenario.Scenario, Expected: scenario.Rules}
		path := filepath.Join(outputDirectory, "config", scenario.Scenario+".yaml")
		_, err := config.LoadConfigFromPath(path)

		validationErrors := validate.Errors(err)
		if len(validationErrors) > 0 {
			for _, validationError := range validationErrors {
				result.Actual = append(result.Actual, validationError.RuleID)
			}
			result.Passed = containsAllRules(result.Actual, scenario.Rules)
		} else if err == nil {
			result.Error = "configuration unexpectedly passed validation"
		} else {
			result.Error = err.Error()
		}

		if result.Passed {
			summary.Passed++
		} else {
			summary.Failed++
		}
		summary.Results = append(summary.Results, result)
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return summary, fmt.Errorf("encode run summary: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDirectory, "results.json"), data, 0644); err != nil {
		return summary, fmt.Errorf("write run summary: %w", err)
	}

	return summary, nil
}

func containsAllRules(actual, expected []validate.RuleID) bool {
	for _, expectedRule := range expected {
		found := false
		for _, actualRule := range actual {
			if actualRule == expectedRule {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
