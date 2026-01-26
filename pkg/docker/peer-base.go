package docker

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"

type PeerBaseNode struct {
	*yaml.Node
}

func NewPeerBase(network string) *PeerBaseNode {
	node := yaml.MappingNode(
		yaml.ScalarNode("peer.base"),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode("hyperledger/fabric-peer:latest"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.MappingNode(yaml.ScalarNode("CORE_VM_ENDPOINT"), yaml.ScalarNode("unix:///host/var/run/docker.sock")),
				yaml.MappingNode(yaml.ScalarNode("CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE"), yaml.ScalarNode(network)),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_LOGGING_SPEC"), yaml.ScalarNode("INFO")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_GOSSIP_USELEADERELECTION"), yaml.ScalarNode("true")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_GOSSIP_ORGLEADER"), yaml.ScalarNode("false")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_GOSSIP_STATE_ENABLED"), yaml.ScalarNode("true")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_PROFILE_ENABLED"), yaml.ScalarNode("true")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_MSPCONFIGPATH"), yaml.ScalarNode("/etc/hyperledger/fabric/msp")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_TLS_ENABLED"), yaml.ScalarNode("true")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_TLS_KEY_FILE"), yaml.ScalarNode("/etc/hyperledger/fabric/tls/server.key")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_TLS_CERT_FILE"), yaml.ScalarNode("/etc/hyperledger/fabric/tls/server.crt")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_TLS_ROOTCERT_FILE"), yaml.ScalarNode("/etc/hyperledger/fabric/tls/ca.crt")),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/etc/hyperledger/fabric"),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("peer node start"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode("./../artifacts/channel:/etc/hyperledger/fabric/channel"),
				yaml.ScalarNode("./../artifacts/chaincode:/etc/hyperledger/fabric/chaincode"),
				yaml.ScalarNode("./../artifacts/crypto-materials:/etc/hyperledger/fabric/crypto-materials"),
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
