package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestInvalidChannelNameFn(t *testing.T) {
	for _, name := range []string{"MYCHANNEL", "my channel", "1channel"} {
		t.Run(name, func(t *testing.T) {
			assertValidationError(t, InvalidChannelNameFn(spec.Channel{Name: name}), RuleChannelNameInvalid, "Invalid Channel Name", "invalid channel name: "+name+" (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)")
		})
	}
	assertNoError(t, InvalidChannelNameFn(spec.Channel{Name: "my-channel.1"}))
}

func TestEmptyChannelNameFn(t *testing.T) {
	assertValidationError(t, EmptyChannelNameFn(spec.Channel{}), RuleChannelNameRequired, "Empty Channel Name", "channel name cannot be empty")
	assertNoError(t, EmptyChannelNameFn(spec.Channel{Name: "mychannel"}))
}

func TestEmptyChannelProfileFn(t *testing.T) {
	channel := spec.Channel{Name: "mychannel"}
	assertValidationError(t, EmptyChannelProfileFn(channel), RuleChannelProfileRequired, "Empty Channel Profile", "channel mychannel must reference a profile")
	channel.Profile = "TwoOrgs"
	assertNoError(t, EmptyChannelProfileFn(channel))
}
