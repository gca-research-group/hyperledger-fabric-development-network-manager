package configtx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/network"
)

type ProfileNode struct {
	*yaml.Node
}

func NewProfile(
	profile config.Profile,
	organizations []config.Organization,
	orderer *OrdererNode,
	applicationDefaults *yaml.Node,
	channelDefaults *yaml.Node,
	applicationOrganizations []*yaml.Node,
	appCapability *yaml.Node,
) *ProfileNode {
	ordererSettings := []*yaml.Node{
		yaml.ScalarNode("<<"),
		yaml.AliasNode(OrdererDefaultsKey, orderer.WithAnchor(OrdererDefaultsKey)),
	}

	if strings.ToLower(profile.Consensus.Type) == strings.ToLower(etcdraftKey) {
		ordererSettings = append(ordererSettings, yaml.ScalarNode(OrdererTypeKey), yaml.ScalarNode(etcdraftKey))
		ordererSettings = append(ordererSettings, BuildRaft(organizations)...)
	} else {
		ordererSettings = append(ordererSettings, yaml.ScalarNode(OrdererTypeKey), yaml.ScalarNode(BFTKey))
		ordererSettings = append(ordererSettings, BuildSmartBFT(organizations)...)
	}

	node := yaml.MappingNode(yaml.ScalarNode(profile.Name),
		yaml.MappingNode(
			yaml.ScalarNode("<<"),
			yaml.AliasNode(ChannelDefaultsKey, channelDefaults),
			yaml.ScalarNode(OrdererKey),
			yaml.MappingNode(
				ordererSettings...,
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

func BuildSmartBFT(organizations []config.Organization) []*yaml.Node {
	var nodes []*yaml.Node

	for _, organization := range organizations {
		for index, orderer := range organization.Orderers {
			mspID := config.ResolveOrdererMSPID(organization)
			nodes = append(nodes,
				yaml.MappingNode(
					yaml.ScalarNode(IDKey),
					yaml.ScalarNode(strconv.Itoa(index+1)),
					yaml.ScalarNode(MSPIDKey),
					yaml.ScalarNode(mspID),
					yaml.ScalarNode(HostKey),
					yaml.ScalarNode(fmt.Sprintf("%s.%s", orderer.Subdomain, organization.Domain)),
					yaml.ScalarNode(PortKey),
					yaml.ScalarNode(fmt.Sprint(orderer.Port)),
					yaml.ScalarNode(ClientTLSCertKey),
					yaml.ScalarNode(fmt.Sprintf(network.ORDERER_TLS_SERVER_CRT, organization.Domain, orderer.Subdomain)),
					yaml.ScalarNode(ServerTLSCertKey),
					yaml.ScalarNode(fmt.Sprintf(network.ORDERER_TLS_SERVER_CRT, organization.Domain, orderer.Subdomain)),
					yaml.ScalarNode(IdentityKey),
					yaml.ScalarNode(fmt.Sprintf(network.ORDERER_MSP_SIGNCERT, organization.Domain, orderer.Subdomain)),
				),
			)
		}
	}

	return []*yaml.Node{
		yaml.ScalarNode(ConsenterMappingKey),
		yaml.SequenceNode(nodes...),
		yaml.ScalarNode(SmartBFTKey),
		yaml.MappingNode(
			yaml.ScalarNode(RequestBatchMaxCountKey), yaml.ScalarNode("100"),
			yaml.ScalarNode(RequestBatchMaxIntervalKey), yaml.ScalarNode("50ms"),
			yaml.ScalarNode(RequestForwardTimeoutKey), yaml.ScalarNode("2s"),
			yaml.ScalarNode(RequestComplainTimeoutKey), yaml.ScalarNode("20s"),
			yaml.ScalarNode(RequestAutoRemoveTimeoutKey), yaml.ScalarNode("3m0s"),
			yaml.ScalarNode(ViewChangeResendIntervalKey), yaml.ScalarNode("5s"),
			yaml.ScalarNode(ViewChangeTimeoutKey), yaml.ScalarNode("20s"),
			yaml.ScalarNode(LeaderHeartbeatTimeoutKey), yaml.ScalarNode("1m0s"),
			yaml.ScalarNode(CollectTimeoutKey), yaml.ScalarNode("1s"),
			yaml.ScalarNode(RequestBatchMaxBytesKey), yaml.ScalarNode("10485760"),
			yaml.ScalarNode(IncomingMessageBufferSizeKey), yaml.ScalarNode("200"),
			yaml.ScalarNode(RequestPoolSizeKey), yaml.ScalarNode("100000"),
			yaml.ScalarNode(LeaderHeartbeatCountKey), yaml.ScalarNode("10"),
		),
	}
}

func BuildRaft(organizations []config.Organization) []*yaml.Node {

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
					yaml.ScalarNode(fmt.Sprintf(network.ORDERER_TLS_SERVER_CRT, organization.Domain, orderer.Subdomain)),
					yaml.ScalarNode(ServerTLSCertKey),
					yaml.ScalarNode(fmt.Sprintf(network.ORDERER_TLS_SERVER_CRT, organization.Domain, orderer.Subdomain)),
				),
			)
		}
	}

	return []*yaml.Node{
		yaml.ScalarNode(BatchTimeoutKey), yaml.ScalarNode(constants.DEFAULT_BATCH_TIMEOUT),
		yaml.ScalarNode(BatchSizeKey), yaml.MappingNode(
			yaml.ScalarNode(MaxMessageCountKey), yaml.ScalarNode(constants.DEFAULT_BATCH_SIZE_MAX_MESSAGE_COUNT),
			yaml.ScalarNode(AbsoluteMaxBytesKey), yaml.ScalarNode(constants.DEFAULT_BATCH_SIZE_ABSOLUTE_MAX_BYTES),
			yaml.ScalarNode(PreferredMaxBytesKey), yaml.ScalarNode(constants.DEFAULT_BATCH_SIZE_PREFERRED_MAX_BYTES),
		),
		yaml.ScalarNode(EtcdRaftKey), yaml.MappingNode(yaml.ScalarNode(ConsentersKey), yaml.SequenceNode(nodes...)),
	}
}

func (pn *ProfileNode) Build() *yaml.Node {
	return pn.Node
}
