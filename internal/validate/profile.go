package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
)

func InvalidConsensusTypeFn(profile spec.Profile) error {
	consensusType := profile.Consensus.Type

	if consensusType != "" && consensusType != "etcdraft" && consensusType != "BFT" {
		return &ValidationError{
			RuleID: RuleConsensusTypeInvalid,
			Rule:   "Invalid Consensus Type",
			Detail: fmt.Sprintf("invalid consensus type for the profile %s", profile.Name),
		}
	}

	return nil
}

func DuplicateProfileNameFn(p spec.Profile, seen map[string]struct{}) error {
	return duplicateValue(p.Name, seen, RuleProfileNameDuplicate, "Duplicate Profile Name", fmt.Sprintf("duplicate profile name: %s", p.Name))
}

func EmptyProfileNameFn(p spec.Profile, i int) error {
	if p.Name == "" {
		return validationError(RuleProfileNameRequired, "Empty Profile Name", fmt.Sprintf("name of the profile index %d is empty", i))
	}
	return nil
}

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
