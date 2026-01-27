package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type CertificateAuthorityNode struct {
	host string
	*yaml.Node
}

func NewCertificateAuthority(host string) *CertificateAuthorityNode {

	node := yaml.MappingNode(
		yaml.ScalarNode(host),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode("hyperledger/fabric-ca:latest"),
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(host),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=%s", host)),
				yaml.ScalarNode("FABRIC_CA_SERVER_TLS_ENABLED=true"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/%s-cert.pem", host)),
				yaml.ScalarNode("FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/priv_sk"),
				yaml.ScalarNode(fmt.Sprintf("FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/%s-cert.pem", host)),
				yaml.ScalarNode("FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/priv_sk"),
			),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("sh -c 'fabric-ca-server start -b admin:adminpw'"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(yaml.ScalarNode("./artifacts/crypto-materials/peerOrganizations/org1.example.com/ca/:/etc/hyperledger/fabric-ca-server-config")),
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
