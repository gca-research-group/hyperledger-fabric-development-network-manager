package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type ToolsNode struct {
	*yaml.Node
}

func NewTools(name string, domain string, corePeerHost string, mspID string, network string) *ToolsNode {
	node := yaml.MappingNode(
		yaml.ScalarNode(fmt.Sprintf("hyperledger-fabric-tools-%s", name)),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger-fabric-tools-%s", name)),
			yaml.ScalarNode("image"),
			yaml.ScalarNode("hyperledger/fabric-tools:latest"),
			yaml.ScalarNode("tty"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("stdin_open"),
			yaml.ScalarNode("open"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.MappingNode(yaml.ScalarNode("GOPATH"), yaml.ScalarNode("/opt/gopath")),
				yaml.MappingNode(yaml.ScalarNode("CORE_VM_ENDPOINT"), yaml.ScalarNode("unix:///host/var/run/docker.sock")),
				yaml.MappingNode(yaml.ScalarNode("FABRIC_LOGGING_SPEC"), yaml.ScalarNode("INFO")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_ID"), yaml.ScalarNode("cli")),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_ADDRESS"), yaml.ScalarNode(corePeerHost)),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_LOCALMSPID"), yaml.ScalarNode(mspID)),
				yaml.MappingNode(yaml.ScalarNode("CORE_PEER_TLS_ENABLED"), yaml.ScalarNode("true")),
				yaml.MappingNode(
					yaml.ScalarNode("CORE_PEER_TLS_CERT_FILE"),
					yaml.ScalarNode(fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/peers/peer0.%s/tls/server.crt", domain, domain)),
				),
				yaml.MappingNode(
					yaml.ScalarNode("CORE_PEER_TLS_KEY_FILE"),
					yaml.ScalarNode(fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/peers/peer0.%s/tls/server.key", domain, domain)),
				),
				yaml.MappingNode(
					yaml.ScalarNode("CORE_PEER_TLS_ROOTCERT_FILE"),
					yaml.ScalarNode(fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/peers/peer0.%s/tls/ca.crt", domain, domain)),
				),
				yaml.MappingNode(
					yaml.ScalarNode("CORE_PEER_MSPCONFIGPATH"),
					yaml.ScalarNode(fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/users/Admin@%s/msp", domain, domain)),
				),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/opt/gopath/src/github.com/hyperledger/fabric/"),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("/bin/bash"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("./artifacts/crypto-materials/peerOrganizations/%s:/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s", domain, domain)),
			),
			yaml.ScalarNode("networks"),
			yaml.SequenceNode(yaml.ScalarNode(network)),
		),
	)
	return &ToolsNode{node}
}

func (tn *ToolsNode) Build() *yaml.Node {
	return tn.Node
}
