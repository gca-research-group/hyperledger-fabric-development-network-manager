package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidOrdererPortFn(orderer spec.Orderer, organizationName string) error {
	if orderer.ExposePort < 0 {
		return &ValidationError{
			RuleID: RuleOrdererPortInvalid,
			Rule:   "Invalid Orderer Port",
			Detail: fmt.Sprintf("expose port of the orderer %s of the organization %s should be greater than zero", orderer.Name, organizationName),
		}
	}

	return nil
}
