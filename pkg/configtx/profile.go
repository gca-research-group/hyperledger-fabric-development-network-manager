package configtx

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type ProfileNode struct {
	*yaml.Node
}

func NewProfile(
	name string,
	ordererDefaults *yaml.Node,
	applicationDefaults *yaml.Node,
	channelDefaults *yaml.Node,
	applicationOrganizations []*yaml.Node,
	appCapability *yaml.Node,
) *ProfileNode {
	node := yaml.MappingNode(yaml.ScalarNode(name),
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(ChannelDefaultsKey, channelDefaults),
			yaml.ScalarNode(OrdererKey),
			yaml.MappingNode(
				yaml.ScalarNode("<<"),
				yaml.AliasNode(OrdererDefaultsKey, ordererDefaults),
			),

			yaml.ScalarNode(ApplicationKey),
			yaml.MappingNode(
				yaml.ScalarNode("<<"),
				yaml.AliasNode(ApplicationDefaultsKey, applicationDefaults),
				yaml.ScalarNode(OrganizationsKey),
				yaml.SequenceNode(applicationOrganizations...),
				yaml.ScalarNode(CapabilitiesKey),
				yaml.MappingNode(
					yaml.ScalarNode("<<"),
					yaml.AliasNode(ApplicationCapabilitiesKey, appCapability),
				),
			),
		),
	)

	return &ProfileNode{node}
}

func (pn *ProfileNode) Build() *yaml.Node {
	return pn.Node
}
