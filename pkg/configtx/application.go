package configtx

type ApplicationNode struct {
	*Node
}

func NewApplication() *ApplicationNode {
	return &ApplicationNode{MappingNode()}
}

func (an *ApplicationNode) WithPolicies() *ApplicationNode {
	an.GetOrCreateValue(PoliciesKey, MappingNode(
		ScalarNode(LifecycleEndorsementKey), NewImplicitMetaPolicy(Policy{Rule: EndorsementKey, Qualifier: MAJORITYKey}),
		ScalarNode(EndorsementKey), NewImplicitMetaPolicy(Policy{Rule: EndorsementKey, Qualifier: MAJORITYKey}),
		ScalarNode(AdminsKey), NewImplicitMetaPolicy(Policy{Rule: AdminsKey, Qualifier: MAJORITYKey}),
		ScalarNode(ReadersKey), NewImplicitMetaPolicy(Policy{Rule: ReadersKey}),
		ScalarNode(WritersKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
	))

	return an
}

func (on *ApplicationNode) WithCapabilities(node *Node) *ApplicationNode {
	on.GetOrCreateValue(CapabilitiesKey,
		MappingNode(
			ScalarNode("<<"),
			AliasNode(ApplicationCapabilitiesKey, node),
		),
	)

	return on
}

func (on *ApplicationNode) WithOrganizations(nodes []*Node) *ApplicationNode {
	on.GetOrCreateValue(OrganizationsKey, SequenceNode(nodes...))
	return on
}
