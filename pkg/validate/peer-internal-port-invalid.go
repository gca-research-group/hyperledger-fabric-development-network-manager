package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidPeerInternalPortFn(p spec.Peer, org string) error {
	if !validOptionalTCPPort(p.Port) {
		return validationError(RulePeerInternalPortInvalid, "Invalid Peer Internal Port", fmt.Sprintf("internal port of peer %s of organization %s must be between 1 and 65535 when set", p.Name, org))
	}
	return nil
}
