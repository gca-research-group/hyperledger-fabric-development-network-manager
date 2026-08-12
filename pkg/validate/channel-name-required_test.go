package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyChannelNameFn(t *testing.T) {
	assertValidationError(t, EmptyChannelNameFn(spec.Channel{}), RuleChannelNameRequired, "Empty Channel Name", "channel name cannot be empty")
	assertNoError(t, EmptyChannelNameFn(spec.Channel{Name: "mychannel"}))
}
