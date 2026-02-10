package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type OrdererNode struct {
	name string
	*yaml.Node
}

func NewOrderer(hostname string, domain string) *OrdererNode {
	ordererDomain := fmt.Sprintf("%s.%s", hostname, domain)

	node := yaml.MappingNode(
		yaml.ScalarNode(ordererDomain),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-orderer:%s", FABRIC_VERSION)),
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(ordererDomain),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/var/hyperledger/orderer"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("ORDERER_GENERAL_LOGLEVEL=INFO"),
				yaml.ScalarNode("ORDERER_GENERAL_LISTENADDRESS=0.0.0.0"),
				yaml.ScalarNode("ORDERER_GENERAL_BOOTSTRAPMETHOD=none"),
				yaml.ScalarNode("ORDERER_GENERAL_LOCALMSPID=OrdererMSP"),
				yaml.ScalarNode("ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_ENABLED=true"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]"),
				yaml.ScalarNode("ORDERER_ADMIN_LISTENADDRESS=0.0.0.0:7053"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_ENABLED=true"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_CLIENTROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]"),
				yaml.ScalarNode("ORDERER_CHANNELPARTICIPATION_ENABLED=true"),
			),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("./%s/certificates/organizations/ordererOrganizations/%s/orderers/%s/msp:/var/hyperledger/orderer/msp", domain, domain, ordererDomain)),
				yaml.ScalarNode(fmt.Sprintf("./%s/certificates/organizations/ordererOrganizations/%s/orderers/%s/tls:/var/hyperledger/orderer/tls", domain, domain, ordererDomain)),
			),
		),
	)

	return &OrdererNode{ordererDomain, node}
}

func (o *OrdererNode) WithNetworks(nodes []*yaml.Node) *OrdererNode {
	node := o.GetValue(o.name)
	node.GetOrCreateValue("networks", yaml.SequenceNode(nodes...))
	return o
}

// TODO:
/* func (o *OrdererNode) WithPort(port int) *OrdererNode {
	node := o.GetValue(o.name)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:7050", port))))
	return o
} */

func (o *OrdererNode) Build() *yaml.Node {
	return o.Node
}
