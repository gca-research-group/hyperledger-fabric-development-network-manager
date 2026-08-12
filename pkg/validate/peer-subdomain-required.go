package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyPeerSubdomainFn(peer spec.Peer, organizationName string) error {
	if peer.Subdomain == "" {
		return &ValidationError{
			RuleID: RulePeerSubdomainRequired,
			Rule:   "Empty Peer Subdomain",
			Detail: fmt.Sprintf("subdomain of the peer %s of the organization %s is undefined", peer.Name, organizationName),
		}
	}

	return nil
}
