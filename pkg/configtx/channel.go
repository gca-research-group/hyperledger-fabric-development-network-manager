package configtx

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"

type ChannelNode struct {
	*yaml.Node
}

func NewChannel() *ChannelNode {
	return &ChannelNode{yaml.MappingNode()}
}

func (ch *ChannelNode) WithPolicies() *ChannelNode {
	ch.GetOrCreateValue(PoliciesKey, yaml.MappingNode(
		yaml.ScalarNode(ReadersKey), NewImplicitMetaPolicy(Policy{Rule: ReadersKey}),
		yaml.ScalarNode(WritersKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
		yaml.ScalarNode(AdminsKey), NewImplicitMetaPolicy(Policy{Rule: AdminsKey, Qualifier: MAJORITYKey}),
	))

	return ch
}

func (on *ChannelNode) WithCapabilities(node *yaml.Node) *ChannelNode {
	on.GetOrCreateValue(CapabilitiesKey,
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(ChannelCapabilitiesKey, node),
		),
	)

	return on
}
