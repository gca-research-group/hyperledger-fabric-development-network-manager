package cryptoconfig

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type PeerOrgNode struct {
	*yaml.Node
}

func NewPeerOrg(organization pkg.Organization) *PeerOrgNode {
	node := yaml.MappingNode(
		yaml.ScalarNode("Name"),
		yaml.ScalarNode(organization.Name),
		yaml.ScalarNode("Domain"),
		yaml.ScalarNode(organization.Domain),
		yaml.ScalarNode("Template"),
		yaml.MappingNode(
			yaml.ScalarNode("Count"), yaml.ScalarNode("1"),
		),
		yaml.ScalarNode("Users"),
		yaml.MappingNode(
			yaml.ScalarNode("Count"), yaml.ScalarNode("1"),
		),
	)

	return &PeerOrgNode{node}
}

func (po *PeerOrgNode) Build() *yaml.Node {
	return po.Node
}
