package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestUndefinedProfileOrganizationFn(t *testing.T) {
	names := map[string]struct{}{"Org1": {}}
	invalid := spec.Profile{Organizations: []string{"Org1", "MissingOrg"}}
	assertValidationError(t, UndefinedProfileOrganizationFn(invalid, names), RuleProfileOrganizationUndefined, "Profile References Undefined Org", "organization not defined: MissingOrg")
	assertNoError(t, UndefinedProfileOrganizationFn(spec.Profile{Organizations: []string{"Org1"}}, names))
}
