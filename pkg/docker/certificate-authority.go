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
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_HOME"), yaml.ScalarNode("/etc/hyperledger/fabric-ca-server")),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_SERVER_CA_NAME"), yaml.ScalarNode(host)),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_SERVER_TLS_ENABLED"), yaml.ScalarNode("true")),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_SERVER_TLS_CERTFILE"), yaml.ScalarNode(fmt.Sprintf("/etc/hyperledger/fabric-ca-server-config/%s-cert.pem", host))),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_SERVER_TLS_KEYFILE"), yaml.ScalarNode("/etc/hyperledger/fabric-ca-server-config/priv_sk")),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_SERVER_CA_CERTFILE"), yaml.ScalarNode(fmt.Sprintf("/etc/hyperledger/fabric-ca-server-config/%s-cert.pem", host))),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_CA_SERVER_CA_KEYFILE"), yaml.ScalarNode("/etc/hyperledger/fabric-ca-server-config/priv_sk")),
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
