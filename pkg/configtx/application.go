package configtx

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"

type ApplicationNode struct {
	*yaml.Node
}

func NewApplication() *ApplicationNode {
	return &ApplicationNode{yaml.MappingNode()}
}

func (an *ApplicationNode) WithPolicies() *ApplicationNode {
	an.GetOrCreateValue(PoliciesKey, yaml.MappingNode(
		yaml.ScalarNode(LifecycleEndorsementKey), NewImplicitMetaPolicy(Policy{Rule: EndorsementKey, Qualifier: MAJORITYKey}),
		yaml.ScalarNode(EndorsementKey), NewImplicitMetaPolicy(Policy{Rule: EndorsementKey, Qualifier: MAJORITYKey}),
		yaml.ScalarNode(AdminsKey), NewImplicitMetaPolicy(Policy{Rule: AdminsKey, Qualifier: MAJORITYKey}),
		yaml.ScalarNode(ReadersKey), NewImplicitMetaPolicy(Policy{Rule: ReadersKey}),
		yaml.ScalarNode(WritersKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
	))

	return an
}

func (on *ApplicationNode) WithCapabilities(node *yaml.Node) *ApplicationNode {
	on.GetOrCreateValue(CapabilitiesKey,
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(ApplicationCapabilitiesKey, node),
		),
	)

	return on
}

func (on *ApplicationNode) WithOrganizations(nodes []*yaml.Node) *ApplicationNode {
	on.GetOrCreateValue(OrganizationsKey, yaml.SequenceNode(nodes...))
	return on
}
