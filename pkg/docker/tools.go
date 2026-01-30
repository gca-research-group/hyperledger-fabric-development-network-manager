package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type ToolsNode struct {
	*yaml.Node
	domain string
}

func NewTools(name string, domain string, corePeerHost string, mspID string, network string) *ToolsNode {
	basePath := "/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials"

	volumes := []*yaml.Node{
		yaml.ScalarNode(fmt.Sprintf("./%s/crypto-config.yml:/opt/gopath/src/github.com/hyperledger/fabric/crypto-config.yml", domain)),
		yaml.ScalarNode("./configtx.yml:/opt/gopath/src/github.com/hyperledger/fabric/configtx.yml"),
		yaml.ScalarNode(fmt.Sprintf("./%s/crypto-materials:/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials", domain)),
		yaml.ScalarNode(fmt.Sprintf("./%s/channel:/opt/gopath/src/github.com/hyperledger/fabric/channel", domain)),
	}

	node := yaml.MappingNode(
		yaml.ScalarNode(fmt.Sprintf("hyperledger-fabric-%s-tools", domain)),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger-fabric-%s-tools", domain)),
			yaml.ScalarNode("image"),
			yaml.ScalarNode("hyperledger/fabric-tools:latest"),
			yaml.ScalarNode("tty"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("stdin_open"),
			yaml.ScalarNode("true"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("GOPATH=/opt/gopath"),
				yaml.ScalarNode("CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock"),
				yaml.ScalarNode("FABRIC_LOGGING_SPEC=INFO"),
				yaml.ScalarNode("CORE_PEER_ID=cli"),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ADDRESS=%s", corePeerHost)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", mspID)),
				yaml.ScalarNode("CORE_PEER_TLS_ENABLED=true"),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=%s/%s/peerOrganizations/peers/peer0.%s/tls/server.crt", basePath, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=%s/%s/peerOrganizations/peers/peer0.%s/tls/server.key", basePath, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=%s/peerOrganizations/%s/peers/peer0.%s/tls/ca.crt", basePath, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%s/peerOrganizations/%s/users/Admin@%s/msp", basePath, domain, domain)),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/opt/gopath/src/github.com/hyperledger/fabric/"),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("/bin/bash"),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(volumes...),
			yaml.ScalarNode("networks"),
			yaml.SequenceNode(yaml.ScalarNode(network)),
		),
	)

	return &ToolsNode{node, domain}
}

func (tn *ToolsNode) WithPeerMSPs(domains []string) *ToolsNode {
	service := tn.GetValue(fmt.Sprintf("hyperledger-fabric-%s-tools", tn.domain))
	volumes := service.GetValue("volumes")

	for _, domain := range domains {
		volume := fmt.Sprintf("./%s/crypto-materials/peerOrganizations/%s/msp:/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/%s/msp", domain, domain, domain)

		var hasVolume bool

		for _, content := range volumes.Content {
			hasVolume = content.Value == volume
		}

		if !hasVolume {
			entry, _ := yaml.ScalarNode(volume).MarshalYAML()
			volumes.Content = append(volumes.Content, entry)
		}
	}

	return tn
}

func (tn *ToolsNode) Build() *yaml.Node {
	return tn.Node
}
