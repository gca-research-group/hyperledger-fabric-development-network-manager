package validate

import (
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func TestUnsupportedApplicationCapabilityFn(t *testing.T) {
	invalid := spec.Config{Capabilities: spec.Capabilities{Application: "V1_4"}}
	assertValidationError(t, UnsupportedApplicationCapabilityFn(invalid), RuleApplicationCapabilityUnsupported, "Unsupported Application Capability", "unsupported application capability: V1_4")
	invalid.Capabilities.Application = "V2_5"
	assertNoError(t, UnsupportedApplicationCapabilityFn(invalid))
}
