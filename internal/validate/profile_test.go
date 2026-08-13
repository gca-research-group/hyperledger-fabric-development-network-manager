package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestInvalidConsensusTypeFn(t *testing.T) {
	profile := spec.Profile{Name: "TwoOrgs", Consensus: spec.Consensus{Type: "raft"}}
	assertValidationError(t, InvalidConsensusTypeFn(profile), RuleConsensusTypeInvalid, "Invalid Consensus Type", "invalid consensus type for the profile TwoOrgs")
	for _, kind := range []string{"", "etcdraft", "BFT"} {
		t.Run(kind, func(t *testing.T) { profile.Consensus.Type = kind; assertNoError(t, InvalidConsensusTypeFn(profile)) })
	}
}

func TestUndefinedProfileOrganizationFn(t *testing.T) {
	names := map[string]struct{}{"Org1": {}}
	invalid := spec.Profile{Organizations: []string{"Org1", "MissingOrg"}}
	assertValidationError(t, UndefinedProfileOrganizationFn(invalid, names), RuleProfileOrganizationUndefined, "Profile References Undefined Org", "organization not defined: MissingOrg")
	assertNoError(t, UndefinedProfileOrganizationFn(spec.Profile{Organizations: []string{"Org1"}}, names))
}

func TestEmptyProfileOrganizationsFn(t *testing.T) {
	profile := spec.Profile{Name: "TwoOrgs"}
	assertValidationError(t, EmptyProfileOrganizationsFn(profile), RuleProfileOrganizationsRequired, "Empty Profile Orgs", "profile TwoOrgs must include at least one organization")
	profile.Organizations = []string{"Org1"}
	assertNoError(t, EmptyProfileOrganizationsFn(profile))
}
