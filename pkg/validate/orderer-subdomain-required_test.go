package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyOrdererSubdomainFn(t *testing.T) {
	orderer := spec.Orderer{Name: "orderer0"}
	assertValidationError(t, EmptyOrdererSubdomainFn(orderer, "Org1"), RuleOrdererSubdomainRequired, "Empty Orderer Subdomain", "subdomain of the orderer orderer0 of the organization Org1 is undefined")
	orderer.Subdomain = "orderer0"
	assertNoError(t, EmptyOrdererSubdomainFn(orderer, "Org1"))
}
