package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidOrdererVersionFn(orderer spec.Orderer, organizationName, channelCapability string) error {
	if err := validateBinary(orderer.Version, spec.MinBinaryVersion[channelCapability]); err != nil {
		return &ValidationError{
			RuleID: RuleOrdererVersionInvalid,
			Rule:   "Invalid Orderer Version",
			Detail: fmt.Sprintf("orderer version of org %s invalid: %v", organizationName, err),
		}
	}

	return nil
}
