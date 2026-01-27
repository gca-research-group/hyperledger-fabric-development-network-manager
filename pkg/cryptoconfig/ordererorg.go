package cryptoconfig

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type OrdererOrgNode struct {
	*yaml.Node
}

func NewOrdererOrg(orderer pkg.Orderer) *OrdererOrgNode {
	node := yaml.MappingNode(
		yaml.ScalarNode("Name"),
		yaml.ScalarNode(orderer.Name),
		yaml.ScalarNode("Domain"),
		yaml.ScalarNode(orderer.Domain),
		yaml.SequenceNode(
			yaml.MappingNode(yaml.ScalarNode("Hostname"), yaml.ScalarNode(orderer.Hostname)),
		),
	)

	return &OrdererOrgNode{node}
}

func (o *OrdererOrgNode) Build() *yaml.Node {
	return o.Node
}
