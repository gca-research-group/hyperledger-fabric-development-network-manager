package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicatePeerNameFn(p spec.Peer, org string, seen map[string]struct{}) error {
	return duplicateValue(p.Name, seen, RulePeerNameDuplicate, "Duplicate Peer Name", fmt.Sprintf("duplicate peer name in organization %s: %s", org, p.Name))
}
