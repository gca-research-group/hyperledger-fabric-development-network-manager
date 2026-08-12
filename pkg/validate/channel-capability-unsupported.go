package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func UnsupportedChannelCapabilityFn(configuration spec.Config) error {
	if _, ok := spec.CapabilityMap[configuration.Capabilities.Channel]; !ok {
		return &ValidationError{
			RuleID: RuleChannelCapabilityUnsupported,
			Rule:   "Unsupported Channel Capability",
			Detail: fmt.Sprintf("unsupported channel capability: %s", configuration.Capabilities.Channel),
		}
	}

	return nil
}
