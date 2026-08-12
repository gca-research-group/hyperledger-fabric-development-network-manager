package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestInvalidConsensusTypeFn(t *testing.T) {
	profile := spec.Profile{Name: "TwoOrgs", Consensus: spec.Consensus{Type: "raft"}}
	assertValidationError(t, InvalidConsensusTypeFn(profile), RuleConsensusTypeInvalid, "Invalid Consensus Type", "invalid consensus type for the profile TwoOrgs")
	for _, kind := range []string{"", "etcdraft", "BFT"} {
		t.Run(kind, func(t *testing.T) { profile.Consensus.Type = kind; assertNoError(t, InvalidConsensusTypeFn(profile)) })
	}
}
