package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestInvalidOrdererVersionFn(t *testing.T) {
	assertValidationError(t, InvalidOrdererVersionFn(spec.Orderer{Version: "2.4.0"}, "Org1", "V2_5"), RuleOrdererVersionInvalid, "Invalid Orderer Version", "orderer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0")
	for _, version := range []string{"", "2.5.0", "3.0.0"} {
		t.Run(version, func(t *testing.T) {
			assertNoError(t, InvalidOrdererVersionFn(spec.Orderer{Version: version}, "Org1", "V2_5"))
		})
	}
}
