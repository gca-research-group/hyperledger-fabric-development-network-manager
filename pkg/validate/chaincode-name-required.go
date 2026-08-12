package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyChaincodeNameFn(chaincode spec.Chaincode, index int) error {
	if chaincode.Name == "" {
		return &ValidationError{
			RuleID: RuleChaincodeNameRequired,
			Rule:   "Empty Chaincode Name",
			Detail: fmt.Sprintf("name of the chaincode %d is empty", index),
		}
	}

	return nil
}
