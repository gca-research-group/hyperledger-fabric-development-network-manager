package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestInvalidOrdererPortFn(t *testing.T) {
	orderer := spec.Orderer{Name: "orderer0", ExposePort: -1}
	assertValidationError(t, InvalidOrdererPortFn(orderer, "Org1"), RuleOrdererPortInvalid, "Invalid Orderer Port", "expose port of the orderer orderer0 of the organization Org1 should be greater than zero")
	orderer.ExposePort = 0
	assertNoError(t, InvalidOrdererPortFn(orderer, "Org1"))
	orderer.ExposePort = 7050
	assertNoError(t, InvalidOrdererPortFn(orderer, "Org1"))
}
