package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestInvalidPeerPortFn(t *testing.T) {
	peer := spec.Peer{Name: "peer0", ExposePort: -1}
	assertValidationError(t, InvalidPeerPortFn(peer, "Org1"), RulePeerPortInvalid, "Invalid Peer Port", "expose port of the peer peer0 of the organization Org1 should be greater than zero")
	peer.ExposePort = 0
	assertNoError(t, InvalidPeerPortFn(peer, "Org1"))
	peer.ExposePort = 7051
	assertNoError(t, InvalidPeerPortFn(peer, "Org1"))
}
