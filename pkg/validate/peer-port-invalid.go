package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidPeerPortFn(peer spec.Peer, organizationName string) error {
	if !validOptionalTCPPort(peer.ExposePort) {
		return &ValidationError{
			RuleID: RulePeerPortInvalid,
			Rule:   "Invalid Peer Port",
			Detail: fmt.Sprintf("expose port of the peer %s of the organization %s must be between 1 and 65535 when set", peer.Name, organizationName),
		}
	}

	return nil
}
