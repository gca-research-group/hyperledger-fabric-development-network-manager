package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyChaincodeVersionFn(t *testing.T) {
	assertValidationError(t, EmptyChaincodeVersionFn(spec.Chaincode{}, 4), RuleChaincodeVersionRequired, "Empty Chaincode Version", "version of the chaincode 4 is empty")
	assertNoError(t, EmptyChaincodeVersionFn(spec.Chaincode{Version: "1.0"}, 0))
}
