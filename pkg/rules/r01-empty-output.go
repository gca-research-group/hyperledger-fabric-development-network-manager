package rules

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/config"

func EmptyOutput(c config.Config) *config.ValidationError {
	if c.Output == "" {
		return &config.ValidationError{
			RuleID: "R01",
			Rule:   "Empty Output",
			Detail: "the output directory cannot be empty",
		}
	}

	return nil
}
