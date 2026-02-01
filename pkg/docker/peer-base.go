package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type PeerBaseNode struct {
	*yaml.Node
}

func NewPeerBase(network string) *PeerBaseNode {
	node := yaml.MappingNode(
		yaml.ScalarNode("peer.base"),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-peer:%s", FABRIC_VERSION)),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock"),
				yaml.ScalarNode("CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=network"),
				yaml.ScalarNode("FABRIC_LOGGING_SPEC=INFO"),
				yaml.ScalarNode("CORE_PEER_GOSSIP_USELEADERELECTION=true"),
				yaml.ScalarNode("CORE_PEER_GOSSIP_ORGLEADER=false"),
				yaml.ScalarNode("CORE_PEER_GOSSIP_STATE_ENABLED=true"),
				yaml.ScalarNode("CORE_PEER_PROFILE_ENABLED=true"),
				yaml.ScalarNode("CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp"),
				yaml.ScalarNode("CORE_PEER_TLS_ENABLED=true"),
				yaml.ScalarNode("CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key"),
				yaml.ScalarNode("CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt"),
				yaml.ScalarNode("CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt"),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/etc/hyperledger/fabric"),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("peer node start"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				// yaml.ScalarNode("./../channel:/etc/hyperledger/fabric/channel"),
				// yaml.ScalarNode("./../chaincode:/etc/hyperledger/fabric/chaincode"),
				// yaml.ScalarNode("./../crypto-materials:/etc/hyperledger/fabric/crypto-materials"),
				yaml.ScalarNode("/var/run/docker.sock:/host/var/run/docker.sock"),
			),
			yaml.ScalarNode("networks"),
			yaml.SequenceNode(yaml.ScalarNode(network)),
		),
	)

	return &PeerBaseNode{node}
}

func (pb *PeerBaseNode) Build() *yaml.Node {
	return pb.Node
}
