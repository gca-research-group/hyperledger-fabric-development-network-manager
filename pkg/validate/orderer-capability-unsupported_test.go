package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestUnsupportedOrdererCapabilityFn(t *testing.T) {
	invalid := spec.Config{Capabilities: spec.Capabilities{Orderer: "V1_4"}}
	assertValidationError(t, UnsupportedOrdererCapabilityFn(invalid), RuleOrdererCapabilityUnsupported, "Unsupported Orderer Capability", "unsupported orderer capability: V1_4")
	invalid.Capabilities.Orderer = "V2_5"
	assertNoError(t, UnsupportedOrdererCapabilityFn(invalid))
}
