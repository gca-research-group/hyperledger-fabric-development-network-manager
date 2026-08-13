package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"regexp"
)

func MultipleBootstrapOrganizationsFn(count int) error {
	if count > 1 {
		return &ValidationError{
			RuleID: RuleBootstrapOrganizationsMultiple,
			Rule:   "Multiple Bootstrap Orgs",
			Detail: "exactly one bootstrap organization must be defined",
		}
	}

	return nil
}

var domainRegex = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+)$`)

func InvalidDomainFn(o spec.Organization) error {
	if o.Domain != "" && !domainRegex.MatchString(o.Domain) {
		return validationError(RuleDomainInvalid, "Invalid Domain", fmt.Sprintf("invalid organization domain: %s", o.Domain))
	}
	return nil
}

func DuplicateOrganizationDomainFn(o spec.Organization, seen map[string]struct{}) error {
	return duplicateValue(o.Domain, seen, RuleOrganizationDomainDuplicate, "Duplicate Organization Domain", fmt.Sprintf("duplicate organization domain: %s", o.Domain))
}

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

func InvalidOrganizationUsersFn(o spec.Organization) error {
	if o.Users < 0 {
		return validationError(RuleOrganizationUsersInvalid, "Invalid Organization Users", fmt.Sprintf("users of organization %s cannot be negative", o.Name))
	}
	return nil
}

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
