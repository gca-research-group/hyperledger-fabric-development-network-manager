package validate

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"

func EmptyNetworkNameFn(c spec.Config) error {
	if c.Network == "" {
		return validationError(RuleNetworkNameRequired, "Empty Network Name", "network name cannot be empty")
	}
	return nil
}
