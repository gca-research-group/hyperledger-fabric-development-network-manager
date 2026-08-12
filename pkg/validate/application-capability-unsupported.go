package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
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
