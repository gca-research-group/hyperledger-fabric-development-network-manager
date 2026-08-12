package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyChaincodePathFn(chaincode spec.Chaincode, index int) error {
	if chaincode.Path == "" {
		return &ValidationError{RuleID: RuleChaincodePathRequired, Rule: "Empty Chaincode Path", Detail: fmt.Sprintf("path of the chaincode %d is empty", index)}
	}

	return nil
}
