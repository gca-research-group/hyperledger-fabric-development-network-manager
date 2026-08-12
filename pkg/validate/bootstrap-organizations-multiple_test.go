package validate

import "testing"

func TestMultipleBootstrapOrganizationsFn(t *testing.T) {
	assertValidationError(t, MultipleBootstrapOrganizationsFn(2), RuleBootstrapOrganizationsMultiple, "Multiple Bootstrap Orgs", "exactly one bootstrap organization must be defined")
	assertNoError(t, MultipleBootstrapOrganizationsFn(0))
	assertNoError(t, MultipleBootstrapOrganizationsFn(1))
}
