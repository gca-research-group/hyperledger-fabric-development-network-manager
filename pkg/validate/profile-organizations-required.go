package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyProfileOrganizationsFn(profile spec.Profile) error {
	if len(profile.Organizations) == 0 {
		return &ValidationError{
			RuleID: RuleProfileOrganizationsRequired,
			Rule:   "Empty Profile Orgs",
			Detail: fmt.Sprintf("profile %s must include at least one organization", profile.Name),
		}
	}

	return nil
}
