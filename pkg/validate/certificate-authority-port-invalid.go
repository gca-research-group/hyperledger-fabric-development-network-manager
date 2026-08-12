package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidCertificateAuthorityPortFn(organization spec.Organization) error {
	if organization.CertificateAuthority.ExposePort < 0 {
		return &ValidationError{
			RuleID: RuleCertificateAuthorityPortInvalid,
			Rule:   "Invalid Certificate Authority Port",
			Detail: fmt.Sprintf("expose port of the certificate authority of the organization %s should be greater than zero", organization.Name),
		}
	}

	return nil
}
