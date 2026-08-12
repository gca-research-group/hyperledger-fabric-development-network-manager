package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyOrganizationDomainFn(t *testing.T) {
	assertValidationError(t, EmptyOrganizationDomainFn(spec.Organization{}, 1), RuleOrganizationDomainRequired, "Empty Organization Domain", "domain of the organization index 1 is undefined")
	assertNoError(t, EmptyOrganizationDomainFn(spec.Organization{Domain: "example.com"}, 0))
}
