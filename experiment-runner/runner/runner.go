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

type Status string

const (
	StatusPassed  Status = "passed"
	StatusPartial Status = "partial"
	StatusFailed  Status = "failed"
)

type Result struct {
	Scenario string            `json:"scenario"`
	Expected []validate.RuleID `json:"expectedRules"`
	Actual   []validate.RuleID `json:"actualRules,omitempty"`
	Missing  []validate.RuleID `json:"missingRules,omitempty"`
	Status   Status            `json:"status"`
	Error    string            `json:"error,omitempty"`
}

type Summary struct {
	Total   int      `json:"total"`
	Passed  int      `json:"passed"`
	Partial int      `json:"partial"`
	Failed  int      `json:"failed"`
	Results []Result `json:"results"`
}

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

func Run(outputDirectory string, scenarios []generator.ScenarioRules) (Summary, error) {
	summary := Summary{Total: len(scenarios), Results: make([]Result, 0, len(scenarios))}

	for _, scenario := range scenarios {
		result := Result{Scenario: scenario.Scenario, Expected: scenario.Rules}
		path := filepath.Join(outputDirectory, "config", scenario.Scenario+".yaml")
		_, err := config.LoadConfigFromPath(path)

		validationErrors := validate.Errors(err)
		if len(validationErrors) > 0 {
			for _, validationError := range validationErrors {
				result.Actual = appendRuleOnce(result.Actual, validationError.RuleID)
			}
			result.Missing = missingRules(result.Actual, scenario.Rules)
			result.Status = classify(len(scenario.Rules), len(result.Missing))
		} else if err == nil {
			result.Status = StatusFailed
			result.Missing = append(result.Missing, scenario.Rules...)
			result.Error = "configuration unexpectedly passed validation"
		} else {
			result.Status = StatusFailed
			result.Missing = append(result.Missing, scenario.Rules...)
			result.Error = err.Error()
		}

		switch result.Status {
		case StatusPassed:
			summary.Passed++
		case StatusPartial:
			summary.Partial++
		case StatusFailed:
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

func missingRules(actual, expected []validate.RuleID) []validate.RuleID {
	var missing []validate.RuleID
	for _, expectedRule := range expected {
		found := false
		for _, actualRule := range actual {
			if actualRule == expectedRule {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, expectedRule)
		}
	}
	return missing
}

func classify(expected, missing int) Status {
	if missing == 0 {
		return StatusPassed
	}
	if missing < expected {
		return StatusPartial
	}
	return StatusFailed
}

func appendRuleOnce(rules []validate.RuleID, ruleID validate.RuleID) []validate.RuleID {
	for _, existing := range rules {
		if existing == ruleID {
			return rules
		}
	}
	return append(rules, ruleID)
}
