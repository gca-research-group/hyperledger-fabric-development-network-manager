package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestUnsupportedChannelCapabilityFn(t *testing.T) {
	invalid := spec.Config{Capabilities: spec.Capabilities{Channel: "V1_4"}}
	assertValidationError(t, UnsupportedChannelCapabilityFn(invalid), RuleChannelCapabilityUnsupported, "Unsupported Channel Capability", "unsupported channel capability: V1_4")
	invalid.Capabilities.Channel = "V2_5"
	assertNoError(t, UnsupportedChannelCapabilityFn(invalid))
}
