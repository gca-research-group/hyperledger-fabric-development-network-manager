package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type OrdererNode struct {
	name string
	*yaml.Node
}

func NewOrderer(name string) *OrdererNode {
	node := yaml.MappingNode(
		yaml.ScalarNode(name),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode("hyperledger/fabric-orderer:2.5"),
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(name),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/var/hyperledger/orderer"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_LOGLEVEL"), yaml.ScalarNode("INFO")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_LISTENADDRESS"), yaml.ScalarNode("0.0.0.0")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_GENESISMETHOD"), yaml.ScalarNode("file")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_GENESISFILE"), yaml.ScalarNode("/var/hyperledger/orderer/orderer.genesis.block")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_LOCALMSPID"), yaml.ScalarNode("OrdererMSP")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_LOCALMSPDIR"), yaml.ScalarNode("/var/hyperledger/orderer/msp")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_TLS_ENABLED"), yaml.ScalarNode("true")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_TLS_PRIVATEKEY"), yaml.ScalarNode("/var/hyperledger/orderer/tls/server.key")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_TLS_CERTIFICATE"), yaml.ScalarNode("/var/hyperledger/orderer/tls/server.crt")),
				yaml.MappingNode(yaml.ScalarNode("ORDERER_GENERAL_TLS_ROOTCAS"), yaml.ScalarNode("[/var/hyperledger/orderer/tls/ca.crt]")),
			),
		),
	)

	return &OrdererNode{name, node}
}

func (o *OrdererNode) WithNetworks(nodes []*yaml.Node) *OrdererNode {
	node := o.GetValue(o.name)
	node.GetOrCreateValue("networks", yaml.SequenceNode(nodes...))
	return o
}

func (o *OrdererNode) WithPort(port int) *OrdererNode {
	node := o.GetValue(o.name)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:7050", port))))
	return o
}

func (o *OrdererNode) Build() *yaml.Node {
	return o.Node
}
