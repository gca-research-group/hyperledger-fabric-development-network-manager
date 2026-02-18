package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type OrdererNode struct {
	*yaml.Node
}

func NewOrderer() *OrdererNode {
	node := yaml.MappingNode(
		yaml.ScalarNode(OrdererTypeKey),
		yaml.ScalarNode(etcdraftKey),
	)

	return &OrdererNode{node}
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

func (on *OrdererNode) WithAddresses(addresses []string) *OrdererNode {
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

func (on *OrdererNode) WithBatchConfig() *OrdererNode {
	on.GetOrCreateValue(BatchTimeoutKey, yaml.ScalarNode(constants.DEFAULT_BATCH_TIMEOUT))
	on.GetOrCreateValue(BatchSizeKey, yaml.MappingNode(
		yaml.ScalarNode(MaxMessageCountKey), yaml.ScalarNode(constants.DEFAULT_BATCH_SIZE_MAX_MESSAGE_COUNT),
		yaml.ScalarNode(AbsoluteMaxBytesKey), yaml.ScalarNode(constants.DEFAULT_BATCH_SIZE_ABSOLUTE_MAX_BYTES),
		yaml.ScalarNode(PreferredMaxBytesKey), yaml.ScalarNode(constants.DEFAULT_BATCH_SIZE_PREFERRED_MAX_BYTES),
	))

	return on
}

func (on *OrdererNode) WithRaftConfig(organizations []config.Organization) *OrdererNode {

	var nodes []*yaml.Node

	for _, organization := range organizations {
		for _, orderer := range organization.Orderers {
			nodes = append(nodes,
				yaml.MappingNode(
					yaml.ScalarNode(HostKey),
					yaml.ScalarNode(fmt.Sprintf("%s.%s", orderer.Subdomain, organization.Domain)),
					yaml.ScalarNode(PortKey),
					yaml.ScalarNode(fmt.Sprint(orderer.Port)),
					yaml.ScalarNode(ClientTLSCertKey),
					yaml.ScalarNode(fmt.Sprintf("./%[1]s/ordererOrganizations/%[1]s/orderers/%[2]s.%[1]s/tls/server.crt", organization.Domain, orderer.Subdomain)),
					yaml.ScalarNode(ServerTLSCertKey),
					yaml.ScalarNode(fmt.Sprintf("./%[1]s/ordererOrganizations/%[1]s/orderers/%[2]s.%[1]s/tls/server.crt", organization.Domain, orderer.Subdomain)),
				),
			)
		}
	}

	on.GetOrCreateValue(EtcdRaftKey,
		yaml.MappingNode(yaml.ScalarNode(ConsentersKey), yaml.SequenceNode(nodes...)),
	)
	return on
}

func (on *OrdererNode) Build() *yaml.Node {
	return on.Node
}
