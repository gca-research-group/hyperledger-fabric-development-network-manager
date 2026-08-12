package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyProfileOrganizationsFn(t *testing.T) {
	profile := spec.Profile{Name: "TwoOrgs"}
	assertValidationError(t, EmptyProfileOrganizationsFn(profile), RuleProfileOrganizationsRequired, "Empty Profile Orgs", "profile TwoOrgs must include at least one organization")
	profile.Organizations = []string{"Org1"}
	assertNoError(t, EmptyProfileOrganizationsFn(profile))
}
