package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestEmptyPeerSubdomainFn(t *testing.T) {
	peer := spec.Peer{Name: "peer0"}
	assertValidationError(t, EmptyPeerSubdomainFn(peer, "Org1"), RulePeerSubdomainRequired, "Empty Peer Subdomain", "subdomain of the peer peer0 of the organization Org1 is undefined")
	peer.Subdomain = "peer0"
	assertNoError(t, EmptyPeerSubdomainFn(peer, "Org1"))
}
