package validate

import (
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func TestConfigReturnsAllValidationErrors(t *testing.T) {
	configuration := spec.Config{
		Output: "invalid:name",
		Capabilities: spec.Capabilities{
			Channel:     "V1_4",
			Application: "V1_4",
			Orderer:     "V1_4",
		},
	}

	rules := make(map[RuleID]bool)
	for _, validationError := range Errors(Config(configuration)) {
		rules[validationError.RuleID] = true
	}
	want := []RuleID{
		RuleOutputDirectoryNameInvalid,
		RuleOrganizationsRequired,
		RuleChannelCapabilityUnsupported,
		RuleApplicationCapabilityUnsupported,
		RuleOrdererCapabilityUnsupported,
		RuleOrdererTopologyRequired,
	}

	for _, ruleID := range want {
		if !rules[ruleID] {
			t.Errorf("missing validation error for rule %s", ruleID)
		}
	}
}
