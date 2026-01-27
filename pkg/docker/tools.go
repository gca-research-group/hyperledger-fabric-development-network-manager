package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type ToolsNode struct {
	*yaml.Node
}

func NewTools(name string, domain string, corePeerHost string, mspID string, network string) *ToolsNode {
	basePath := "/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations"

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
				yaml.ScalarNode("GOPATH=/opt/gopath"),
				yaml.ScalarNode("CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock"),
				yaml.ScalarNode("FABRIC_LOGGING_SPEC=INFO"),
				yaml.ScalarNode("CORE_PEER_ID=cli"),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ADDRESS=%s", corePeerHost)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", mspID)),
				yaml.ScalarNode("CORE_PEER_TLS_ENABLED=true"),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=%s/%s/peers/peer0.%s/tls/server.crt", basePath, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=%s/%s/peers/peer0.%s/tls/server.key", basePath, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=%s/%s/peers/peer0.%s/tls/ca.crt", basePath, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s/%s/users/Admin@%s/msp", basePath, domain, domain)),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/opt/gopath/src/github.com/hyperledger/fabric/"),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("/bin/bash"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("./artifacts/crypto-materials/peerOrganizations/%s:%s/%s", domain, basePath, domain)),
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
