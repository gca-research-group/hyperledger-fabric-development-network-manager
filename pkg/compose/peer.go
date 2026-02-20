package compose

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type PeerNode struct {
	*yaml.Node
	peer   config.Peer
	domain string
}

func NewPeer(
	mspID string,
	peer config.Peer,
	currentOrganization config.Organization,
	corePeerGossipBootstrap string,
	network string,
	organizations []config.Organization,
) *PeerNode {

	domain := currentOrganization.Domain
	peerDomain := ResolvePeerDomain(peer.Subdomain, domain)
	peerPort := ResolvePeerPort(peer.Port)

	node := yaml.MappingNode(
		yaml.ScalarNode(peerDomain),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(peerDomain),
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-peer:%s", ResolvePeerVersion(currentOrganization.Version.Peer))),
			yaml.ScalarNode("extends"),
			yaml.MappingNode(
				yaml.ScalarNode("file"),
				yaml.ScalarNode("./peer.base.yml"),
				yaml.ScalarNode("service"),
				yaml.ScalarNode("peer.base"),
			),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", mspID)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ID=%s", peerDomain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ADDRESS=%s:%d", peerDomain, peerPort)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_LISTENADDRESS=0.0.0.0:%d", peerPort)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s", corePeerGossipBootstrap)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:%d", peerDomain, peerPort)),
				yaml.ScalarNode("CORE_LEDGER_STATE_STATEDATABASE=CouchDB"),
				yaml.ScalarNode(fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb.%s:5984", peerDomain)),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin"),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw"),
			),
		),
	)

	return &PeerNode{node, peer, domain}
}

func (pn *PeerNode) ExposePort() *PeerNode {
	if pn.peer.ExposePort == 0 {
		return pn
	}

	domain := pn.domain
	peerDomain := ResolvePeerDomain(pn.peer.Subdomain, domain)

	node := pn.GetValue(peerDomain)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:%d", pn.peer.ExposePort, ResolvePeerPort(pn.peer.Port)))))

	return pn
}

func (pn *PeerNode) WithVolumes() *PeerNode {
	domain := pn.domain
	peerDomain := ResolvePeerDomain(pn.peer.Subdomain, domain)

	peerHostDir := fmt.Sprintf("./%[1]s/certificate-authority/organizations/peerOrganizations/%[1]s/peers/%[2]s", domain, peerDomain)
	peerContainerDir := "/etc/hyperledger/fabric"

	volumes := []*yaml.Node{
		yaml.ScalarNode(fmt.Sprintf("%s/msp:%s/msp", peerHostDir, peerContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("%s/tls:%s/tls", peerHostDir, peerContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("./%s/peers/%s/peer:/var/hyperledger/production", domain, pn.peer.Subdomain)),
	}

	node := pn.GetValue(peerDomain)
	node.GetOrCreateValue("volumes", yaml.SequenceNode(volumes...))

	return pn
}

func (pn *PeerNode) Build() *yaml.Node {
	return pn.Node
}
