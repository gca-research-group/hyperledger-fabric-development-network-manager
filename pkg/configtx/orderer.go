package configtx

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type OrdererNode struct {
	*yaml.Node
}

func NewOrderer(capabilities config.Capabilities) *OrdererNode {
	return &OrdererNode{yaml.MappingNode()}
}

func (on *OrdererNode) WithCapabilities(node *yaml.Node) *OrdererNode {
	on.GetOrCreateValue(CapabilitiesKey,
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(OrdererCapabilitiesKey, node),
		),
	)

	return on
}

func (on *OrdererNode) WithAddresses(addresses []string, capabilities config.Capabilities) *OrdererNode {
	if capabilities.Channel == "V3_0" {
		return on
	}

	var nodes []*yaml.Node

	for _, address := range addresses {
		nodes = append(nodes, yaml.ScalarNode(address))
	}

	on.GetOrCreateValue(AddressesKey, yaml.SequenceNode(nodes...))
	return on
}

func (on *OrdererNode) WithPolicies() *OrdererNode {
	on.GetOrCreateValue(PoliciesKey, yaml.MappingNode(
		yaml.ScalarNode(ReadersKey), NewImplicitMetaPolicy(Policy{Rule: ReadersKey}),
		yaml.ScalarNode(WritersKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
		yaml.ScalarNode(AdminsKey), NewImplicitMetaPolicy(Policy{Rule: AdminsKey, Qualifier: MAJORITYKey}),
		yaml.ScalarNode(BlockValidationKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
	))

	return on
}

func (on *OrdererNode) WithOrganizations(nodes []*yaml.Node) *OrdererNode {
	on.GetOrCreateValue(OrganizationsKey, yaml.SequenceNode(nodes...))

	return on
}

func (on *OrdererNode) Build() *yaml.Node {
	return on.Node
}
