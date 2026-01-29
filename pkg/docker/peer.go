package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type PeerNode struct {
	*yaml.Node
}

func NewPeer(
	mspID string,
	host string,
	domain string,
	corePeerGossipBootstrap string,
	network string,
) *PeerNode {
	node := yaml.MappingNode(
		yaml.ScalarNode(host),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(host),
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
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ID=%s", host)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", host)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s", corePeerGossipBootstrap)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:7051", host)),
				yaml.ScalarNode("CORE_LEDGER_STATE_STATEDATABASE=CouchDB"),
				yaml.ScalarNode(fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb.%s:5984", host)),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin"),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=adminpw"),
			),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("./crypto-materials/peerOrganizations/%s/peers/%s/msp:/etc/hyperledger/fabric/msp", domain, host)),
				yaml.ScalarNode(fmt.Sprintf("./crypto-materials/peerOrganizations/%s/peers/%s/tls:/etc/hyperledger/fabric/tls", domain, host)),
			),
		),
		yaml.ScalarNode(fmt.Sprintf("couchdb.%s", host)),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(fmt.Sprintf("couchdb.%s", host)),
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
