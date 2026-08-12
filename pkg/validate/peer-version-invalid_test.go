package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestInvalidPeerVersionFn(t *testing.T) {
	assertValidationError(t, InvalidPeerVersionFn(spec.Peer{Version: "2.4.0"}, "Org1", "V2_5"), RulePeerVersionInvalid, "Invalid Peer Version", "peer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0")
	for _, version := range []string{"", "2.5.0", "3.0.0"} {
		t.Run(version, func(t *testing.T) {
			assertNoError(t, InvalidPeerVersionFn(spec.Peer{Version: version}, "Org1", "V2_5"))
		})
	}
}
