package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyOrganizationNameFn(t *testing.T) {
	assertValidationError(t, EmptyOrganizationNameFn(spec.Organization{}, 2), RuleOrganizationNameRequired, "Empty Organization Name", "name of the organization index 2 is undefined")
	assertNoError(t, EmptyOrganizationNameFn(spec.Organization{Name: "Org1"}, 0))
}
