package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type OrganizationNode struct {
	*yaml.Node
}

func BuildMSPID(name string) string {
	return fmt.Sprintf("%sMSP", name)
}

func NewApplicationOrganization(name string, domain string, mspID string, ordererAddresses []string) *OrganizationNode {
	ordererEndpoints := []*yaml.Node{}

	for _, address := range ordererAddresses {
		ordererEndpoints = append(ordererEndpoints, yaml.ScalarNode(address))
	}

	node := yaml.MappingNode(
		yaml.ScalarNode(NameKey),
		yaml.ScalarNode(name),
		yaml.ScalarNode(IDKey),
		yaml.ScalarNode(mspID),
		yaml.ScalarNode(MSPDirKey),
		yaml.ScalarNode(fmt.Sprintf("./%s/peerOrganizations/%s/msp", domain, domain)),
		yaml.ScalarNode("OrdererEndpoints"),
		yaml.SequenceNode(ordererEndpoints...),
	).WithAnchor(name).WithTag("!!map")

	return &OrganizationNode{node}
}

func NewOrdererOrganization(name string, domain string, mspID string) *OrganizationNode {
	node := yaml.MappingNode(
		yaml.ScalarNode(NameKey),
		yaml.ScalarNode(name),
		yaml.ScalarNode(IDKey),
		yaml.ScalarNode(mspID),
		yaml.ScalarNode(MSPDirKey),
		yaml.ScalarNode(fmt.Sprintf("./%s/ordererOrganizations/%s/msp", domain, domain)),
	).WithAnchor(name).WithTag("!!map")

	return &OrganizationNode{node}
}

func (on *OrganizationNode) WithAnchorPeer(anchorPeer pkg.AnchorPeer) *OrganizationNode {
	if anchorPeer.Host == "" {
		return on
	}

	peer := on.GetOrCreateValue(AnchorPeersKey, yaml.SequenceNode())

	entry, _ := yaml.MappingNode(
		yaml.ScalarNode(HostKey),
		yaml.ScalarNode(anchorPeer.Host),
		yaml.ScalarNode(PortKey),
		yaml.ScalarNode(fmt.Sprint(anchorPeer.Port)),
	).MarshalYAML()

	peer.Content = append(peer.Content, entry)

	return on
}

func (on *OrganizationNode) WithDefaultApplicationPolicies(mspID string) *OrganizationNode {
	policies := yaml.MappingNode(
		yaml.ScalarNode(ReadersKey), NewMemberPolicy(mspID),
		yaml.ScalarNode(WritersKey), NewMemberPolicy(mspID),
		yaml.ScalarNode(AdminsKey), NewAdminPolicy(mspID),
		yaml.ScalarNode(EndorsementKey), NewPeerPolicy(mspID),
	)

	on.GetOrCreateValue(PoliciesKey, policies)
	return on
}

func (on *OrganizationNode) WithDefaultOrdererPolicies(mspID string) *OrganizationNode {
	policies := yaml.MappingNode(
		yaml.ScalarNode(ReadersKey), NewMemberPolicy(mspID),
		yaml.ScalarNode(WritersKey), NewMemberPolicy(mspID),
		yaml.ScalarNode(AdminsKey), NewAdminPolicy(mspID),
	)

	on.GetOrCreateValue(PoliciesKey, policies)
	return on
}

func (on OrganizationNode) Build() *yaml.Node {
	return on.Node
}
