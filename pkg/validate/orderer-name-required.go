package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyOrdererNameFn(orderer spec.Orderer, index int, organizationName string) error {
	if orderer.Name == "" {
		return &ValidationError{
			RuleID: RuleOrdererNameRequired,
			Rule:   "Empty Orderer Name",
			Detail: fmt.Sprintf("name of the orderer index %d of the organization %s is undefined", index, organizationName),
		}
	}

	return nil
}
