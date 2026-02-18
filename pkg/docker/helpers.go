package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
)

func resolvePeerPort(port int) int {
	if port == 0 {
		return constants.DEFAULT_PEER_PORT
	}

	return port
}

func resolveOrdererPort(port int) int {
	if port == 0 {
		return constants.DEFAULT_ORDERER_PORT
	}

	return port
}

func resolvePeerVersion(version string) string {
	if version == "" {
		return constants.DEFAULT_FABRIC_VERSION
	}

	return version
}

func resolvePeerDomain(subdomain string, domain string) string {
	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func resolveOrdererDomain(subdomain string, domain string) string {
	return fmt.Sprintf("%s.%s", subdomain, domain)
}

func resolveOrdererVersion(version string) string {
	if version == "" {
		return constants.DEFAULT_FABRIC_VERSION
	}

	return version
}

func resolveCertificateAuthorityDomain(domain string) string {
	return fmt.Sprintf("ca.%s", domain)
}

func resolveCertificateAuthorityVersion(version string) string {
	if version == "" {
		return "latest"
	}

	return version
}
