package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
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
