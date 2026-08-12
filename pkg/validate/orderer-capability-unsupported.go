package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

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
