package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicateOrganizationNameFn(organization spec.Organization, names map[string]struct{}) error {
	if _, exists := names[organization.Name]; exists {
		return &ValidationError{
			RuleID: RuleOrganizationNameDuplicate,
			Rule:   "Duplicate Org Name",
			Detail: fmt.Sprintf("duplicate organization name: %s", organization.Name),
		}
	}

	return nil
}
