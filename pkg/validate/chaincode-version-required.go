package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyChaincodeVersionFn(chaincode spec.Chaincode, index int) error {
	if chaincode.Version == "" {
		return &ValidationError{
			RuleID: RuleChaincodeVersionRequired,
			Rule:   "Empty Chaincode Version",
			Detail: fmt.Sprintf("version of the chaincode %d is empty", index),
		}
	}

	return nil
}
