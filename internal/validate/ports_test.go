package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestInvalidCertificateAuthorityPortFn(t *testing.T) {
	org := spec.Organization{Name: "Org1", CertificateAuthority: spec.CertificateAuthority{ExposePort: -1}}
	assertValidationError(t, InvalidCertificateAuthorityPortFn(org), RuleCertificateAuthorityPortInvalid, "Invalid Certificate Authority Port", "expose port of the certificate authority of the organization Org1 must be between 1 and 65535 when set")
	org.CertificateAuthority.ExposePort = 0
	assertNoError(t, InvalidCertificateAuthorityPortFn(org))
	org.CertificateAuthority.ExposePort = 7054
	assertNoError(t, InvalidCertificateAuthorityPortFn(org))
}

func TestExposedPortConflictFn(t *testing.T) {
	ports := make(map[int]portOwner)
	peer := portOwner{ownerType: "peer", name: "peer0"}
	assertNoError(t, ExposedPortConflictFn(7051, peer, ports))
	if got := ports[7051]; got != peer {
		t.Fatalf("expected port owner %#v, got %#v", peer, got)
	}
	assertValidationError(t, ExposedPortConflictFn(7051, portOwner{ownerType: "orderer", name: "orderer0"}, ports), RuleExposedPortConflict, "Exposed Port Conflict", "Port 7051 is assigned to both peer 'peer0' and orderer 'orderer0'.")
	assertNoError(t, ExposedPortConflictFn(0, portOwner{ownerType: "peer", name: "ignored"}, ports))
	if _, exists := ports[0]; exists {
		t.Fatal("zero port must not be recorded")
	}
}

func TestInvalidOrdererPortFn(t *testing.T) {
	orderer := spec.Orderer{Name: "orderer0", ExposePort: -1}
	assertValidationError(t, InvalidOrdererPortFn(orderer, "Org1"), RuleOrdererPortInvalid, "Invalid Orderer Port", "expose port of the orderer orderer0 of the organization Org1 must be between 1 and 65535 when set")
	orderer.ExposePort = 0
	assertNoError(t, InvalidOrdererPortFn(orderer, "Org1"))
	orderer.ExposePort = 7050
	assertNoError(t, InvalidOrdererPortFn(orderer, "Org1"))
}

func TestInvalidPeerPortFn(t *testing.T) {
	peer := spec.Peer{Name: "peer0", ExposePort: -1}
	assertValidationError(t, InvalidPeerPortFn(peer, "Org1"), RulePeerPortInvalid, "Invalid Peer Port", "expose port of the peer peer0 of the organization Org1 must be between 1 and 65535 when set")
	peer.ExposePort = 0
	assertNoError(t, InvalidPeerPortFn(peer, "Org1"))
	peer.ExposePort = 7051
	assertNoError(t, InvalidPeerPortFn(peer, "Org1"))
}

func TestTCPPortRange(t *testing.T) {
	for _, port := range []int{0, 1, 65535} {
		assertNoError(t, InvalidPeerInternalPortFn(spec.Peer{Port: port}, "Org1"))
	}
	for _, port := range []int{-1, 65536} {
		checks := []error{
			InvalidCertificateAuthorityPortFn(spec.Organization{Name: "Org1", CertificateAuthority: spec.CertificateAuthority{ExposePort: port}}),
			InvalidPeerPortFn(spec.Peer{ExposePort: port}, "Org1"),
			InvalidPeerInternalPortFn(spec.Peer{Port: port}, "Org1"),
			InvalidOrdererPortFn(spec.Orderer{ExposePort: port}, "Org1"),
			InvalidOrdererInternalPortFn(spec.Orderer{Port: port}, "Org1"),
		}
		for _, err := range checks {
			if err == nil {
				t.Errorf("TCP port %d must be invalid", port)
			}
		}
	}
}
