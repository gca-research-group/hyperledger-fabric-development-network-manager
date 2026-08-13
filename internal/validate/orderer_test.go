package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestEmptyOrdererNameFn(t *testing.T) {
	assertValidationError(t, EmptyOrdererNameFn(spec.Orderer{}, 2, "Org1"), RuleOrdererNameRequired, "Empty Orderer Name", "name of the orderer index 2 of the organization Org1 is undefined")
	assertNoError(t, EmptyOrdererNameFn(spec.Orderer{Name: "orderer0"}, 0, "Org1"))
}

func TestEmptyOrdererSubdomainFn(t *testing.T) {
	orderer := spec.Orderer{Name: "orderer0"}
	assertValidationError(t, EmptyOrdererSubdomainFn(orderer, "Org1"), RuleOrdererSubdomainRequired, "Empty Orderer Subdomain", "subdomain of the orderer orderer0 of the organization Org1 is undefined")
	orderer.Subdomain = "orderer0"
	assertNoError(t, EmptyOrdererSubdomainFn(orderer, "Org1"))
}

func TestInvalidOrdererVersionFn(t *testing.T) {
	assertValidationError(t, InvalidOrdererVersionFn(spec.Orderer{Version: "2.4.0"}, "Org1", "V2_5"), RuleOrdererVersionInvalid, "Invalid Orderer Version", "orderer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0")
	for _, version := range []string{"", "2.5.0", "3.0.0"} {
		t.Run(version, func(t *testing.T) {
			assertNoError(t, InvalidOrdererVersionFn(spec.Orderer{Version: version}, "Org1", "V2_5"))
		})
	}
}
