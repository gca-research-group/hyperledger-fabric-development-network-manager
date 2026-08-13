package validate

import (
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

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
