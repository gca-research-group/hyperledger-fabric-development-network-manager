package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyChaincodeNameFn(t *testing.T) {
	assertValidationError(t, EmptyChaincodeNameFn(spec.Chaincode{}, 4), RuleChaincodeNameRequired, "Empty Chaincode Name", "name of the chaincode 4 is empty")
	assertNoError(t, EmptyChaincodeNameFn(spec.Chaincode{Name: "asset"}, 0))
}
