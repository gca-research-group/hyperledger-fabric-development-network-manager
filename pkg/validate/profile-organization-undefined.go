package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func UndefinedProfileOrganizationFn(profile spec.Profile, organizationNames map[string]struct{}) error {
	for _, name := range profile.Organizations {
		if _, ok := organizationNames[name]; !ok {
			return &ValidationError{
				RuleID: RuleProfileOrganizationUndefined,
				Rule:   "Profile References Undefined Org",
				Detail: fmt.Sprintf("organization not defined: %s", name),
			}
		}
	}

	return nil
}
