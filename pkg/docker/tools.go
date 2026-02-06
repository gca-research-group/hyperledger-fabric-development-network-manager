package docker

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type ToolsNode struct {
	*yaml.Node
	name string
}

func NewTools(currentOrganization pkg.Organization, organizations []pkg.Organization, network string) *ToolsNode {
	name := currentOrganization.Name
	domain := currentOrganization.Domain
	corePeerHost := fmt.Sprintf("peer0.%s:7051", currentOrganization.Domain)
	mspID := fmt.Sprintf("%sMSP", currentOrganization.Name)

	volumes := []*yaml.Node{
		yaml.ScalarNode("./configtx.yml:/opt/gopath/src/github.com/hyperledger/fabric/configtx.yml"),
		yaml.ScalarNode(fmt.Sprintf("./%s/channel:/opt/gopath/src/github.com/hyperledger/fabric/channel", domain)),

		yaml.ScalarNode(fmt.Sprintf("./%s/certificates/crypto-material:/opt/gopath/src/github.com/hyperledger/fabric/%s", domain, domain)),
	}

	for _, organization := range organizations {
		if organization.Domain == currentOrganization.Domain {
			continue
		}

		var l, r string

		if len(organization.Orderers) > 0 {
			for _, orderer := range organization.Orderers {

				l = fmt.Sprintf("./%s/certificates/crypto-material/ordererOrganizations/%s/orderers/%s.%s/msp/cacerts", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/ordererOrganizations/%s/orderers/%s.%s/msp/cacerts", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

				l = fmt.Sprintf("./%s/certificates/crypto-material/ordererOrganizations/%s/orderers/%s.%s/msp/signcerts", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/ordererOrganizations/%s/orderers/%s.%s/msp/signcerts", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

				l = fmt.Sprintf("./%s/certificates/crypto-material/ordererOrganizations/%s/orderers/%s.%s/msp/tlscacerts", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/ordererOrganizations/%s/orderers/%s.%s/msp/tlscacerts", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

				l = fmt.Sprintf("./%s/certificates/crypto-material/ordererOrganizations/%s/orderers/%s.%s/tls/ca.crt", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/ordererOrganizations/%s/orderers/%s.%s/tls/ca.crt", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

				l = fmt.Sprintf("./%s/certificates/crypto-material/ordererOrganizations/%s/orderers/%s.%s/tls/server.crt", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/ordererOrganizations/%s/orderers/%s.%s/tls/server.crt", organization.Domain, organization.Domain, orderer.Hostname, organization.Domain)
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))
			}
		}

		peers := 1

		if organization.Peers > 0 {
			peers = organization.Peers
		}

		l = fmt.Sprintf("./%s/certificates/crypto-material/peerOrganizations/%s/msp/cacerts", organization.Domain, organization.Domain)
		r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/msp/cacerts", organization.Domain, organization.Domain)
		volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

		for i := 0; i < peers; i++ {

			l = fmt.Sprintf("./%s/certificates/crypto-material/peerOrganizations/%s/peers/peer%d.%s/msp/signcerts", organization.Domain, organization.Domain, i, organization.Domain)
			r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/peers/peer%d.%s/msp/signcerts", organization.Domain, organization.Domain, i, organization.Domain)
			volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

			l = fmt.Sprintf("./%s/certificates/crypto-material/peerOrganizations/%s/peers/peer%d.%s/tls/ca.crt", organization.Domain, organization.Domain, i, organization.Domain)
			r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/peers/peer%d.%s/tls/ca.crt", organization.Domain, organization.Domain, i, organization.Domain)
			volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))

			l = fmt.Sprintf("./%s/certificates/crypto-material/peerOrganizations/%s/peers/peer%d.%s/tls/server.crt", organization.Domain, organization.Domain, i, organization.Domain)
			r = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/peers/peer%d.%s/tls/server.crt", organization.Domain, organization.Domain, i, organization.Domain)
			volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s:%s", l, r)))
		}
	}

	node := yaml.MappingNode(
		yaml.ScalarNode(fmt.Sprintf("hyperledger-fabric-tools-%s", strings.ToLower(name))),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger-fabric-tools-%s", strings.ToLower(name))),
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-tools:%s", FABRIC_VERSION)),
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
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/peers/peer0.%s/tls/server.crt", domain, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/peers/peer0.%s/tls/server.key", domain, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/peers/peer0.%s/tls/ca.crt", domain, domain, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/%s/peerOrganizations/%s/users/Admin@%s/msp", domain, domain, domain)),
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

	return &ToolsNode{node, name}
}

func (tn *ToolsNode) Build() *yaml.Node {
	return tn.Node
}
