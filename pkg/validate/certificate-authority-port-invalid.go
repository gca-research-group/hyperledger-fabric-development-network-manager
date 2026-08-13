package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidCertificateAuthorityPortFn(organization spec.Organization) error {
	if !validOptionalTCPPort(organization.CertificateAuthority.ExposePort) {
		return &ValidationError{
			RuleID: RuleCertificateAuthorityPortInvalid,
			Rule:   "Invalid Certificate Authority Port",
			Detail: fmt.Sprintf("expose port of the certificate authority of the organization %s must be between 1 and 65535 when set", organization.Name),
		}
	}

	return nil
}
