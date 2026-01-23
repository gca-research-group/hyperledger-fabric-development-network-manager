package configtx

type ChannelNode struct {
	*Node
}

func NewChannel() *ChannelNode {
	return &ChannelNode{MappingNode()}
}

func (ch *ChannelNode) WithPolicies() *ChannelNode {
	ch.GetOrCreateValue(PoliciesKey, MappingNode(
		ScalarNode(ReadersKey), NewImplicitMetaPolicy(Policy{Rule: ReadersKey}),
		ScalarNode(WritersKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
		ScalarNode(AdminsKey), NewImplicitMetaPolicy(Policy{Rule: AdminsKey, Qualifier: MAJORITYKey}),
	))

	return ch
}

func (on *ChannelNode) WithCapabilities(node *Node) *ChannelNode {
	on.GetOrCreateValue(CapabilitiesKey,
		MappingNode(
			ScalarNode("<<"),
			AliasNode(ChannelCapabilitiesKey, node),
		),
	)

	return on
}
