package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyPeerNameFn(peer spec.Peer, index int, organizationName string) error {
	if peer.Name == "" {
		return &ValidationError{
			RuleID: RulePeerNameRequired,
			Rule:   "Empty Peer Name",
			Detail: fmt.Sprintf("name of the peer index %d of the organization %s is undefined", index, organizationName),
		}
	}

	return nil
}
