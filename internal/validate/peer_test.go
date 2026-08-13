package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestEmptyPeerNameFn(t *testing.T) {
	assertValidationError(t, EmptyPeerNameFn(spec.Peer{}, 3, "Org1"), RulePeerNameRequired, "Empty Peer Name", "name of the peer index 3 of the organization Org1 is undefined")
	assertNoError(t, EmptyPeerNameFn(spec.Peer{Name: "peer0"}, 0, "Org1"))
}

func TestEmptyPeerSubdomainFn(t *testing.T) {
	peer := spec.Peer{Name: "peer0"}
	assertValidationError(t, EmptyPeerSubdomainFn(peer, "Org1"), RulePeerSubdomainRequired, "Empty Peer Subdomain", "subdomain of the peer peer0 of the organization Org1 is undefined")
	peer.Subdomain = "peer0"
	assertNoError(t, EmptyPeerSubdomainFn(peer, "Org1"))
}

func TestInvalidPeerVersionFn(t *testing.T) {
	assertValidationError(t, InvalidPeerVersionFn(spec.Peer{Version: "2.4.0"}, "Org1", "V2_5"), RulePeerVersionInvalid, "Invalid Peer Version", "peer version of org Org1 invalid: version 2.4.0 is lower than required 2.5.0")
	for _, version := range []string{"", "2.5.0", "3.0.0"} {
		t.Run(version, func(t *testing.T) {
			assertNoError(t, InvalidPeerVersionFn(spec.Peer{Version: version}, "Org1", "V2_5"))
		})
	}
}
