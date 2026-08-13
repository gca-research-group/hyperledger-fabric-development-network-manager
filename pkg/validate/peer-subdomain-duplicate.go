package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicatePeerSubdomainFn(p spec.Peer, org string, seen map[string]struct{}) error {
	return duplicateValue(p.Subdomain, seen, RulePeerSubdomainDuplicate, "Duplicate Peer Subdomain", fmt.Sprintf("duplicate peer subdomain in organization %s: %s", org, p.Subdomain))
}
