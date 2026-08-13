package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
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

func ExposedPortConflictFn(port int, owner portOwner, exposedPorts map[int]portOwner) error {
	if port <= 0 {
		return nil
	}

	if existingOwner, exists := exposedPorts[port]; exists {
		return &ValidationError{
			RuleID: RuleExposedPortConflict,
			Rule:   "Exposed Port Conflict",
			Detail: fmt.Sprintf("Port %d is assigned to both %s and %s.", port, existingOwner, owner),
		}
	}

	exposedPorts[port] = owner

	return nil
}

func InvalidOrdererInternalPortFn(o spec.Orderer, org string) error {
	if !validOptionalTCPPort(o.Port) {
		return validationError(RuleOrdererInternalPortInvalid, "Invalid Orderer Internal Port", fmt.Sprintf("internal port of orderer %s of organization %s must be between 1 and 65535 when set", o.Name, org))
	}
	return nil
}

func InvalidOrdererPortFn(orderer spec.Orderer, organizationName string) error {
	if !validOptionalTCPPort(orderer.ExposePort) {
		return &ValidationError{
			RuleID: RuleOrdererPortInvalid,
			Rule:   "Invalid Orderer Port",
			Detail: fmt.Sprintf("expose port of the orderer %s of the organization %s must be between 1 and 65535 when set", orderer.Name, organizationName),
		}
	}

	return nil
}

func InvalidPeerInternalPortFn(p spec.Peer, org string) error {
	if !validOptionalTCPPort(p.Port) {
		return validationError(RulePeerInternalPortInvalid, "Invalid Peer Internal Port", fmt.Sprintf("internal port of peer %s of organization %s must be between 1 and 65535 when set", p.Name, org))
	}
	return nil
}

func InvalidPeerPortFn(peer spec.Peer, organizationName string) error {
	if !validOptionalTCPPort(peer.ExposePort) {
		return &ValidationError{
			RuleID: RulePeerPortInvalid,
			Rule:   "Invalid Peer Port",
			Detail: fmt.Sprintf("expose port of the peer %s of the organization %s must be between 1 and 65535 when set", peer.Name, organizationName),
		}
	}

	return nil
}
