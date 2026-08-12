package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyPeerNameFn(t *testing.T) {
	assertValidationError(t, EmptyPeerNameFn(spec.Peer{}, 3, "Org1"), RulePeerNameRequired, "Empty Peer Name", "name of the peer index 3 of the organization Org1 is undefined")
	assertNoError(t, EmptyPeerNameFn(spec.Peer{Name: "peer0"}, 0, "Org1"))
}
