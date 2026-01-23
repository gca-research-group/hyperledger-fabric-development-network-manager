package configtx

type ProfileNode struct {
	*Node
}

func NewDefaultProfiles(
	ordererDefaults *Node,
	applicationDefaults *Node,
	channelDefaults *Node,
	ordererOrganizations []*Node,
	applicationOrganizations []*Node,
) *Node {
	return MappingNode(
		ScalarNode(OrdererGenesisProfileKey),
		MappingNode(
			ScalarNode("<<"),
			AliasNode(ChannelDefaultsKey, channelDefaults),
			ScalarNode(OrdererKey),
			MappingNode(
				ScalarNode("<<"),
				AliasNode(OrdererDefaultsKey, ordererDefaults),
				ScalarNode(OrganizationsKey),
				SequenceNode(ordererOrganizations...),
			),
			ScalarNode(ConsortiumsKey),
			MappingNode(
				ScalarNode(DefaultConsortiumKey),
				MappingNode(
					ScalarNode(OrganizationsKey),
					SequenceNode(applicationOrganizations...),
				),
			),
		),
		ScalarNode(SampleProfileKey),
		MappingNode(
			ScalarNode("<<"),
			AliasNode(ChannelDefaultsKey, channelDefaults),
			ScalarNode(ConsortiumKey),
			ScalarNode(DefaultConsortiumKey),
			ScalarNode(ApplicationKey),
			MappingNode(
				ScalarNode("<<"),
				AliasNode(ApplicationDefaultsKey, applicationDefaults),
				ScalarNode(OrganizationsKey),
				SequenceNode(applicationOrganizations...),
			),
		),
	)
}

/* MultiChannel:
<<: *ChannelDefaults
Consortium: MultiConsortium
Application:
  <<: *ApplicationDefaults
  Organizations:
    - *Org1
    - *Org2
    - *Org3
*/

func (pn *ProfileNode) Build() *Node {
	return pn.Node
}
