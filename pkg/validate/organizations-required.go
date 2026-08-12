package validate

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"

func NoOrganizationFn(configuration spec.Config) error {
	if len(configuration.Organizations) == 0 {
		return &ValidationError{
			RuleID: RuleOrganizationsRequired,
			Rule:   "No Organization",
			Detail: "at least one organization must be defined",
		}
	}

	return nil
}
