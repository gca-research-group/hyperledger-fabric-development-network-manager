package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestDuplicateOrganizationNameFn(t *testing.T) {
	names := map[string]struct{}{"Org1": {}}
	assertValidationError(t, DuplicateOrganizationNameFn(spec.Organization{Name: "Org1"}, names), RuleOrganizationNameDuplicate, "Duplicate Org Name", "duplicate organization name: Org1")
	assertNoError(t, DuplicateOrganizationNameFn(spec.Organization{Name: "Org2"}, names))
}
