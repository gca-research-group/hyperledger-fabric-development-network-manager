package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestEmptyChaincodeNameFn(t *testing.T) {
	assertValidationError(t, EmptyChaincodeNameFn(spec.Chaincode{}, 4), RuleChaincodeNameRequired, "Empty Chaincode Name", "name of the chaincode 4 is empty")
	assertNoError(t, EmptyChaincodeNameFn(spec.Chaincode{Name: "asset"}, 0))
}

func TestEmptyChaincodePathFn(t *testing.T) {
	assertValidationError(t, EmptyChaincodePathFn(spec.Chaincode{}, 4), RuleChaincodePathRequired, "Empty Chaincode Path", "path of the chaincode 4 is empty")
	assertNoError(t, EmptyChaincodePathFn(spec.Chaincode{Path: "chaincode/asset"}, 0))
}

func TestEmptyChaincodeVersionFn(t *testing.T) {
	assertValidationError(t, EmptyChaincodeVersionFn(spec.Chaincode{}, 4), RuleChaincodeVersionRequired, "Empty Chaincode Version", "version of the chaincode 4 is empty")
	assertNoError(t, EmptyChaincodeVersionFn(spec.Chaincode{Version: "1.0"}, 0))
}
