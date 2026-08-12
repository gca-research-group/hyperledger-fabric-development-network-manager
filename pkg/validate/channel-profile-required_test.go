package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyChannelProfileFn(t *testing.T) {
	channel := spec.Channel{Name: "mychannel"}
	assertValidationError(t, EmptyChannelProfileFn(channel), RuleChannelProfileRequired, "Empty Channel Profile", "channel mychannel must reference a profile")
	channel.Profile = "TwoOrgs"
	assertNoError(t, EmptyChannelProfileFn(channel))
}
