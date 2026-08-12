package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyOrdererSubdomainFn(orderer spec.Orderer, organizationName string) error {
	if orderer.Subdomain == "" {
		return &ValidationError{
			RuleID: RuleOrdererSubdomainRequired,
			Rule:   "Empty Orderer Subdomain",
			Detail: fmt.Sprintf("subdomain of the orderer %s of the organization %s is undefined", orderer.Name, organizationName),
		}
	}

	return nil
}
