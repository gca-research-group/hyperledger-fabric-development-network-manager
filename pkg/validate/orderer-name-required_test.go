package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyOrdererNameFn(t *testing.T) {
	assertValidationError(t, EmptyOrdererNameFn(spec.Orderer{}, 2, "Org1"), RuleOrdererNameRequired, "Empty Orderer Name", "name of the orderer index 2 of the organization Org1 is undefined")
	assertNoError(t, EmptyOrdererNameFn(spec.Orderer{Name: "orderer0"}, 0, "Org1"))
}
