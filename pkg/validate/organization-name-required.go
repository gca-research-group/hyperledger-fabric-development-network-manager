package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyOrganizationNameFn(organization spec.Organization, index int) error {
	if organization.Name == "" {
		return &ValidationError{
			RuleID: RuleOrganizationNameRequired,
			Rule:   "Empty Organization Name",
			Detail: fmt.Sprintf("name of the organization index %d is undefined", index),
		}
	}

	return nil
}
