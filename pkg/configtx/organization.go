package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"gopkg.in/yaml.v3"
)

type OrganizationNode struct {
	*Node
}

func BuildMSPID(name string) string {
	return fmt.Sprintf("%sMSP", name)
}

func NewApplicationOrganization(name string, domain string, mspID string) *OrganizationNode {
	node := MappingNode(
		ScalarNode(NameKey),
		ScalarNode(name),
		ScalarNode(IDKey),
		ScalarNode(mspID),
		ScalarNode(MSPDirKey),
		ScalarNode(fmt.Sprintf("./crypto-materials/peerOrganizations/%s/msp", domain)),
	).WithAnchor(name).WithTag("!!map")

	return &OrganizationNode{node}
}

func NewOrdererOrganization(name string, domain string, mspID string) *OrganizationNode {
	node := MappingNode(
		ScalarNode(NameKey),
		ScalarNode(name),
		ScalarNode(IDKey),
		ScalarNode(mspID),
		ScalarNode(MSPDirKey),
		ScalarNode(fmt.Sprintf("./crypto-materials/ordererOrganizations/%s/msp", domain)),
	).WithAnchor(name).WithTag("!!map")

	return &OrganizationNode{node}
}

func (on *OrganizationNode) WithAnchorPeer(anchorPeer pkg.AnchorPeer) *OrganizationNode {
	if anchorPeer.Host == "" {
		return on
	}

	peer := on.GetOrCreateValue(AnchorPeersKey, SequenceNode())

	entry := MappingNode(
		ScalarNode(HostKey),
		ScalarNode(anchorPeer.Host),
		ScalarNode(PortKey),
		ScalarNode(fmt.Sprint(anchorPeer.Port)),
	)

	peer.Content = append(peer.Content, (*yaml.Node)(entry))

	return on
}

func (on *OrganizationNode) WithDefaultApplicationPolicies(mspID string) *OrganizationNode {
	policies := MappingNode(
		ScalarNode(ReadersKey), NewPeerPolicy(mspID),
		ScalarNode(WritersKey), NewPeerPolicy(mspID),
		ScalarNode(AdminsKey), NewAdminPolicy(mspID),
		ScalarNode(EndorsementKey), NewPeerPolicy(mspID),
	)

	on.GetOrCreateValue(PoliciesKey, policies)
	return on
}

func (on *OrganizationNode) WithDefaultOrdererPolicies(mspID string) *OrganizationNode {
	policies := MappingNode(
		ScalarNode(ReadersKey), NewMemberPolicy(mspID),
		ScalarNode(WritersKey), NewMemberPolicy(mspID),
		ScalarNode(AdminsKey), NewAdminPolicy(mspID),
	)

	on.GetOrCreateValue(PoliciesKey, policies)
	return on
}

func (on OrganizationNode) Build() *Node {
	return on.Node
}
