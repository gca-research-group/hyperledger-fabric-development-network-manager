package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type PeerNode struct {
	*yaml.Node
}

func NewPeer(
	mspID string,
	peerDomain string,
	domain string,
	corePeerGossipBootstrap string,
	network string,
	organizations []config.Organization,
) *PeerNode {

	peerHostDir := fmt.Sprintf("./%[1]s/certificates/organizations/peerOrganizations/%[1]s/peers/%[2]s", domain, peerDomain)
	peerContainerDir := "/etc/hyperledger/fabric"

	volumes := []*yaml.Node{
		yaml.ScalarNode(fmt.Sprintf("%s/msp:%s/msp", peerHostDir, peerContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("%s/tls:%s/tls", peerHostDir, peerContainerDir)),
	}

	node := yaml.MappingNode(
		yaml.ScalarNode(peerDomain),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(peerDomain),
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
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", peerDomain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s", corePeerGossipBootstrap)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:7051", peerDomain)),
				yaml.ScalarNode("CORE_LEDGER_STATE_STATEDATABASE=CouchDB"),
				yaml.ScalarNode(fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb.%s:5984", peerDomain)),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin"),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw"),
			),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(volumes...),
		),
		yaml.ScalarNode(fmt.Sprintf("couchdb.%s", peerDomain)),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(fmt.Sprintf("couchdb.%s", peerDomain)),
			yaml.ScalarNode("image"),
			yaml.ScalarNode("couchdb:latest"),
			yaml.ScalarNode("environment"),
			yaml.MappingNode(
				yaml.ScalarNode("COUCHDB_USER"),
				yaml.ScalarNode("admin"),
				yaml.ScalarNode("COUCHDB_PASSWORD"),
				yaml.ScalarNode("adminpw"),
			),
			yaml.ScalarNode("networks"),
			yaml.SequenceNode(yaml.ScalarNode(network)),
		),
	)

	return &PeerNode{node}
}

func (np *PeerNode) Build() *yaml.Node {
	return np.Node
}
