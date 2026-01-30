package configtx

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type ProfileNode struct {
	*yaml.Node
}

func NewDefaultProfile(
	ordererDefaults *yaml.Node,
	channelDefaults *yaml.Node,
	ordererOrganizations []*yaml.Node,
	applicationOrganizations []*yaml.Node,
) *ProfileNode {
	node := yaml.MappingNode(
		yaml.ScalarNode(OrdererGenesisProfileKey),
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(ChannelDefaultsKey, channelDefaults),
			yaml.ScalarNode(OrdererKey),
			yaml.MappingNode(
				yaml.ScalarNode("<<"),
				yaml.AliasNode(OrdererDefaultsKey, ordererDefaults),
				yaml.ScalarNode(OrganizationsKey),
				yaml.SequenceNode(ordererOrganizations...),
			),
			yaml.ScalarNode(ConsortiumsKey),
			yaml.MappingNode(
				yaml.ScalarNode(DefaultConsortiumKey),
				yaml.MappingNode(
					yaml.ScalarNode(OrganizationsKey),
					yaml.SequenceNode(applicationOrganizations...),
				),
			),
		),
	)

	return &ProfileNode{node}
}

func NewProfile(
	name string,
	ordererDefaults *yaml.Node,
	applicationDefaults *yaml.Node,
	channelDefaults *yaml.Node,
	ordererOrganizations []*yaml.Node,
	applicationOrganizations []*yaml.Node,
) *ProfileNode {
	node := yaml.MappingNode(yaml.ScalarNode(name),
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(ChannelDefaultsKey, channelDefaults),
			yaml.ScalarNode(ConsortiumKey),
			yaml.ScalarNode(DefaultConsortiumKey),
			yaml.ScalarNode(ApplicationKey),
			yaml.MappingNode(
				yaml.ScalarNode("<<"),
				yaml.AliasNode(ApplicationDefaultsKey, applicationDefaults),
				yaml.ScalarNode(OrganizationsKey),
				yaml.SequenceNode(applicationOrganizations...),
			),
		),
	)

	return &ProfileNode{node}
}

func (pn *ProfileNode) Build() *yaml.Node {
	return pn.Node
}
