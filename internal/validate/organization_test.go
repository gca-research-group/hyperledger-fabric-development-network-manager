package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestMultipleBootstrapOrganizationsFn(t *testing.T) {
	assertValidationError(t, MultipleBootstrapOrganizationsFn(2), RuleBootstrapOrganizationsMultiple, "Multiple Bootstrap Orgs", "exactly one bootstrap organization must be defined")
	assertNoError(t, MultipleBootstrapOrganizationsFn(0))
	assertNoError(t, MultipleBootstrapOrganizationsFn(1))
}

func TestEmptyOrganizationDomainFn(t *testing.T) {
	assertValidationError(t, EmptyOrganizationDomainFn(spec.Organization{}, 1), RuleOrganizationDomainRequired, "Empty Organization Domain", "domain of the organization index 1 is undefined")
	assertNoError(t, EmptyOrganizationDomainFn(spec.Organization{Domain: "example.com"}, 0))
}

func TestDuplicateOrganizationNameFn(t *testing.T) {
	names := map[string]struct{}{"Org1": {}}
	assertValidationError(t, DuplicateOrganizationNameFn(spec.Organization{Name: "Org1"}, names), RuleOrganizationNameDuplicate, "Duplicate Org Name", "duplicate organization name: Org1")
	assertNoError(t, DuplicateOrganizationNameFn(spec.Organization{Name: "Org2"}, names))
}

func TestEmptyOrganizationNameFn(t *testing.T) {
	assertValidationError(t, EmptyOrganizationNameFn(spec.Organization{}, 2), RuleOrganizationNameRequired, "Empty Organization Name", "name of the organization index 2 is undefined")
	assertNoError(t, EmptyOrganizationNameFn(spec.Organization{Name: "Org1"}, 0))
}

func TestNoOrganizationFn(t *testing.T) {
	assertValidationError(t, NoOrganizationFn(spec.Config{}), RuleOrganizationsRequired, "No Organization", "at least one organization must be defined")
	assertNoError(t, NoOrganizationFn(spec.Config{Organizations: []spec.Organization{{Name: "Org1"}}}))
}
