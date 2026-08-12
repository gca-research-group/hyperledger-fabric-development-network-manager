package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyOrganizationDomainFn(organization spec.Organization, index int) error {
	if organization.Domain == "" {
		return &ValidationError{
			RuleID: RuleOrganizationDomainRequired,
			Rule:   "Empty Organization Domain",
			Detail: fmt.Sprintf("domain of the organization index %d is undefined", index),
		}
	}

	return nil
}
