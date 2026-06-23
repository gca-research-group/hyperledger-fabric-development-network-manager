package rules

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/config"

func NoOrganizations(c config.Config) *config.ValidationError {
	if len(c.Organizations) == 0 {
		return &config.ValidationError{
			RuleID: "R02",
			Rule:   "No Organizations",
			Detail: "at least one organization must be defined",
		}
	}

	return nil
}
