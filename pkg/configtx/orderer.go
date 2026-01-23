package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

type OrdererNode struct {
	*Node
}

func NewOrderer() *OrdererNode {
	node := MappingNode(
		ScalarNode(OrdererTypeKey),
		ScalarNode(etcdraftKey),
	)

	return &OrdererNode{node}
}

func (on *OrdererNode) WithCapabilities(node *Node) *OrdererNode {
	on.GetOrCreateValue(CapabilitiesKey,
		MappingNode(
			ScalarNode("<<"),
			AliasNode(OrdererCapabilitiesKey, node),
		),
	)

	return on
}

func (on *OrdererNode) WithAddresses(addresses []string) *OrdererNode {
	var nodes []*Node

	for _, address := range addresses {
		nodes = append(nodes, ScalarNode(address))
	}

	on.GetOrCreateValue(AddressesKey, SequenceNode(nodes...))
	return on
}

func (on *OrdererNode) WithPolicies() *OrdererNode {
	on.GetOrCreateValue(PoliciesKey, MappingNode(
		ScalarNode(ReadersKey), NewImplicitMetaPolicy(Policy{Rule: ReadersKey}),
		ScalarNode(WritersKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
		ScalarNode(AdminsKey), NewImplicitMetaPolicy(Policy{Rule: AdminsKey, Qualifier: MAJORITYKey}),
		ScalarNode(BlockValidationKey), NewImplicitMetaPolicy(Policy{Rule: WritersKey}),
	))
	return on
}

func (on *OrdererNode) WithOrganizations(nodes []*Node) *OrdererNode {
	on.GetOrCreateValue(OrganizationsKey, SequenceNode(nodes...))

	return on
}

func (on *OrdererNode) WithBatchConfig() *OrdererNode {
	on.GetOrCreateValue(BatchTimeoutKey, ScalarNode("2s"))
	on.GetOrCreateValue(BatchSizeKey, MappingNode(
		ScalarNode(MaxMessageCountKey), ScalarNode("10"),
		ScalarNode(AbsoluteMaxBytesKey), ScalarNode("99 MB"),
		ScalarNode(PreferredMaxBytesKey), ScalarNode("512 KB"),
	))

	return on
}

func (on *OrdererNode) WithRaftConfig(orderers []pkg.Orderer) *OrdererNode {

	var nodes []*Node

	for _, orderer := range orderers {
		for _, address := range orderer.Addresses {
			nodes = append(nodes,
				MappingNode(
					ScalarNode(HostKey),
					ScalarNode(address.Host),
					ScalarNode(PortKey),
					ScalarNode(fmt.Sprint(address.Port)),
					ScalarNode(ClientTLSCertKey),
					ScalarNode(fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/%s/tls/server.crt", orderer.Domain, address.Host)),
					ScalarNode(ServerTLSCertKey),
					ScalarNode(fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/%s/tls/server.crt", orderer.Domain, address.Host)),
				),
			)
		}
	}

	on.GetOrCreateValue(EtcdRaftKey,
		MappingNode(ScalarNode(ConsentersKey), SequenceNode(nodes...)),
	)
	return on
}

func (on *OrdererNode) Build() *Node {
	return on.Node
}
