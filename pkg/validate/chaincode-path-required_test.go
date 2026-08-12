package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyChaincodePathFn(t *testing.T) {
	assertValidationError(t, EmptyChaincodePathFn(spec.Chaincode{}, 4), RuleChaincodePathRequired, "Empty Chaincode Path", "path of the chaincode 4 is empty")
	assertNoError(t, EmptyChaincodePathFn(spec.Chaincode{Path: "chaincode/asset"}, 0))
}
