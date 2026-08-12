package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
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
