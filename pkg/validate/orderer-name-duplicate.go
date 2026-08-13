package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicateOrdererNameFn(o spec.Orderer, org string, seen map[string]struct{}) error {
	return duplicateValue(o.Name, seen, RuleOrdererNameDuplicate, "Duplicate Orderer Name", fmt.Sprintf("duplicate orderer name in organization %s: %s", org, o.Name))
}
