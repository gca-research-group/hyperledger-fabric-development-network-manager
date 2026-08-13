package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
)

func UnsupportedApplicationCapabilityFn(configuration spec.Config) error {
	if _, ok := spec.CapabilityMap[configuration.Capabilities.Application]; !ok {
		return &ValidationError{
			RuleID: RuleApplicationCapabilityUnsupported,
			Rule:   "Unsupported Application Capability",
			Detail: fmt.Sprintf("unsupported application capability: %s", configuration.Capabilities.Application),
		}
	}

	return nil
}

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

func UnsupportedOrdererCapabilityFn(configuration spec.Config) error {
	if _, ok := spec.CapabilityMap[configuration.Capabilities.Orderer]; !ok {
		return &ValidationError{
			RuleID: RuleOrdererCapabilityUnsupported,
			Rule:   "Unsupported Orderer Capability",
			Detail: fmt.Sprintf("unsupported orderer capability: %s", configuration.Capabilities.Orderer),
		}
	}

	return nil
}
