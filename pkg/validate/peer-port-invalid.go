package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidPeerPortFn(peer spec.Peer, organizationName string) error {
	if peer.ExposePort < 0 {
		return &ValidationError{
			RuleID: RulePeerPortInvalid,
			Rule:   "Invalid Peer Port",
			Detail: fmt.Sprintf("expose port of the peer %s of the organization %s should be greater than zero", peer.Name, organizationName),
		}
	}

	return nil
}
