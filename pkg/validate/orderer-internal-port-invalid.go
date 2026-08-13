package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidOrdererInternalPortFn(o spec.Orderer, org string) error {
	if !validOptionalTCPPort(o.Port) {
		return validationError(RuleOrdererInternalPortInvalid, "Invalid Orderer Internal Port", fmt.Sprintf("internal port of orderer %s of organization %s must be between 1 and 65535 when set", o.Name, org))
	}
	return nil
}
