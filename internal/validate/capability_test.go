package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestUnsupportedApplicationCapabilityFn(t *testing.T) {
	invalid := spec.Config{Capabilities: spec.Capabilities{Application: "V1_4"}}
	assertValidationError(t, UnsupportedApplicationCapabilityFn(invalid), RuleApplicationCapabilityUnsupported, "Unsupported Application Capability", "unsupported application capability: V1_4")
	invalid.Capabilities.Application = "V2_5"
	assertNoError(t, UnsupportedApplicationCapabilityFn(invalid))
}

func TestUnsupportedChannelCapabilityFn(t *testing.T) {
	invalid := spec.Config{Capabilities: spec.Capabilities{Channel: "V1_4"}}
	assertValidationError(t, UnsupportedChannelCapabilityFn(invalid), RuleChannelCapabilityUnsupported, "Unsupported Channel Capability", "unsupported channel capability: V1_4")
	invalid.Capabilities.Channel = "V2_5"
	assertNoError(t, UnsupportedChannelCapabilityFn(invalid))
}

func TestUnsupportedOrdererCapabilityFn(t *testing.T) {
	invalid := spec.Config{Capabilities: spec.Capabilities{Orderer: "V1_4"}}
	assertValidationError(t, UnsupportedOrdererCapabilityFn(invalid), RuleOrdererCapabilityUnsupported, "Unsupported Orderer Capability", "unsupported orderer capability: V1_4")
	invalid.Capabilities.Orderer = "V2_5"
	assertNoError(t, UnsupportedOrdererCapabilityFn(invalid))
}
