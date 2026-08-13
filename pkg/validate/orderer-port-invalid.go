package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidOrdererPortFn(orderer spec.Orderer, organizationName string) error {
	if !validOptionalTCPPort(orderer.ExposePort) {
		return &ValidationError{
			RuleID: RuleOrdererPortInvalid,
			Rule:   "Invalid Orderer Port",
			Detail: fmt.Sprintf("expose port of the orderer %s of the organization %s must be between 1 and 65535 when set", orderer.Name, organizationName),
		}
	}

	return nil
}
