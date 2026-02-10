package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type CertificateAuthorityNode struct {
	host string
	*yaml.Node
}

func NewCertificateAuthority(organization pkg.Organization) *CertificateAuthorityNode {

	host := fmt.Sprintf("ca.%s", organization.Domain)

	node := yaml.MappingNode(
		yaml.ScalarNode(host),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode("hyperledger/fabric-ca:latest"),
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(host),
			yaml.ScalarNode("tty"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("stdin_open"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=%s", host)),
				yaml.ScalarNode("FABRIC_CA_SERVER_TLS_ENABLED=true"),
				yaml.ScalarNode("FABRIC_CA_SERVER_PORT=7054"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CSR_CN=%s", host)),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CSR_HOSTS=localhost,%s", host)),
				// yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CSR_NAMES=C=US,ST=North Carolina,L=Durham,O=%s,OU=Fabric", organization.Name)),
				// yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/%s-cert.pem", host)),
				// yaml.ScalarNode("FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/priv_sk"),
				// yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/%s-cert.pem", host)),
				// yaml.ScalarNode("FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/priv_sk"),
			),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("sh -c 'fabric-ca-server start -b admin:adminpw'"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("./%s/certificates/fabric-ca-server:/etc/hyperledger/fabric-ca-server", organization.Domain)),
				yaml.ScalarNode(fmt.Sprintf("./%s/certificates/organizations:/etc/hyperledger/organizations", organization.Domain)),
			),
		),
	)

	return &CertificateAuthorityNode{host, node}
}

func (ca *CertificateAuthorityNode) WithNetworks(nodes []*yaml.Node) *CertificateAuthorityNode {
	node := ca.GetValue(ca.host)
	node.GetOrCreateValue("networks", yaml.SequenceNode(nodes...))
	return ca
}

func (ca *CertificateAuthorityNode) WithPort(port int) *CertificateAuthorityNode {
	node := ca.GetValue(ca.host)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:7050", port))))
	return ca
}

func (ca *CertificateAuthorityNode) Build() *yaml.Node {
	return ca.Node
}
