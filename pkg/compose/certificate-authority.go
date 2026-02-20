package compose

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type CertificateAuthorityNode struct {
	*yaml.Node
	certificateAuthority config.CertificateAuthority
	domain               string
}

func NewCertificateAuthority(organization config.Organization) *CertificateAuthorityNode {

	domain := organization.Domain
	certificateAuthority := organization.CertificateAuthority
	certificateAuthorityDomain := ResolveCertificateAuthorityDomain(domain)

	version := ResolveCertificateAuthorityVersion(organization.Version.CertificateAuthority)

	node := yaml.MappingNode(
		yaml.ScalarNode(certificateAuthorityDomain),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-ca:%s", version)),
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(certificateAuthorityDomain),
			yaml.ScalarNode("tty"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("stdin_open"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=%s", certificateAuthorityDomain)),
				yaml.ScalarNode("FABRIC_CA_SERVER_TLS_ENABLED=true"),
				yaml.ScalarNode("FABRIC_CA_SERVER_PORT=7054"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CSR_CN=%s", certificateAuthorityDomain)),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CSR_HOSTS=localhost,%s", certificateAuthorityDomain)),
			),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("sh -c 'fabric-ca-server start -b admin:adminpw'"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("./%s/certificate-authority/fabric-ca-server:/etc/hyperledger/fabric-ca-server", domain)),
				yaml.ScalarNode(fmt.Sprintf("./%s/certificate-authority/organizations:/etc/hyperledger/organizations", domain)),
			),
		),
	)

	return &CertificateAuthorityNode{node, certificateAuthority, domain}
}

func (ca *CertificateAuthorityNode) WithNetworks(nodes []*yaml.Node) *CertificateAuthorityNode {
	node := ca.GetValue(ResolveCertificateAuthorityDomain(ca.domain))
	node.GetOrCreateValue("networks", yaml.SequenceNode(nodes...))
	return ca
}

func (ca *CertificateAuthorityNode) ExposePort() *CertificateAuthorityNode {
	if ca.certificateAuthority.ExposePort == 0 {
		return ca
	}

	certificateAuthorityDomain := ResolveCertificateAuthorityDomain(ca.domain)

	node := ca.GetValue(certificateAuthorityDomain)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:%d", ca.certificateAuthority.ExposePort, constants.DEFAULT_CERTIFICATE_AUTHORITY_PORT))))

	return ca
}

func (ca *CertificateAuthorityNode) Build() *yaml.Node {
	return ca.Node
}
