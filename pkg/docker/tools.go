package docker

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/constants"
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
		yaml.ScalarNode(fmt.Sprintf("./configtx.yml:%s/configtx.yml", constants.DEFAULT_FABRIC_DIRECTORY)),
		yaml.ScalarNode(fmt.Sprintf("./%s/channel:%s/channel", domain, constants.DEFAULT_FABRIC_DIRECTORY)),

		yaml.ScalarNode(fmt.Sprintf("./%[1]s/certificates/organizations:%[2]s/%[1]s", domain, constants.DEFAULT_FABRIC_DIRECTORY)),
	}

	for _, organization := range organizations {
		if organization.Domain == currentOrganization.Domain {
			continue
		}

		if len(organization.Orderers) > 0 {
			for _, orderer := range organization.Orderers {

				ordererHostDir := fmt.Sprintf("./%[1]s/certificates/organizations/ordererOrganizations/%[1]s/orderers/%[2]s.%[1]s", organization.Domain, orderer.Hostname)
				ordererContainerDir := fmt.Sprintf("%[1]s/%[2]s/ordererOrganizations/%[2]s/orderers/%[3]s.%[2]s", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain, orderer.Hostname)

				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/msp/cacerts:%s/msp/cacerts", ordererHostDir, ordererContainerDir)))
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/msp/signcerts:%s/msp/signcerts", ordererHostDir, ordererContainerDir)))
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/msp/tlscacerts:%s/msp/tlscacerts", ordererHostDir, ordererContainerDir)))

				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/tls/ca.crt:%s/tls/ca.crt", ordererHostDir, ordererContainerDir)))
				volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/tls/server.crt:%s/tls/server.crt", ordererHostDir, ordererContainerDir)))
			}
		}

		peerHostDir := fmt.Sprintf("./%[1]s/certificates/organizations/peerOrganizations/%[1]s", organization.Domain)
		peerContainerDir := fmt.Sprintf("%[1]s/%[2]s/peerOrganizations/%[2]s", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain)

		volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/msp/cacerts:%s/msp/cacerts", peerHostDir, peerContainerDir)))
		volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/msp/config.yaml:%s/msp/config.yaml", peerHostDir, peerContainerDir)))

		for i := 0; i < organization.Peers; i++ {
			peerHostDir := fmt.Sprintf("./%[1]s/certificates/organizations/peerOrganizations/%[1]s/peers/peer%[2]d.%[1]s", organization.Domain, i)
			peerContainerDir := fmt.Sprintf("%[1]s/%[2]s/peerOrganizations/%[2]s/peers/peer%[3]d.%[2]s", constants.DEFAULT_FABRIC_DIRECTORY, organization.Domain, i)

			volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/msp/signcerts:%s/msp/signcerts", peerHostDir, peerContainerDir)))

			volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/tls/ca.crt:%s/tls/ca.crt", peerHostDir, peerContainerDir)))
			volumes = append(volumes, yaml.ScalarNode(fmt.Sprintf("%s/tls/server.crt:%s/tls/server.crt", peerHostDir, peerContainerDir)))
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
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=%[1]s/%[2]s/peerOrganizations/%[2]s/peers/peer0.%[2]s/tls/server.crt", constants.DEFAULT_FABRIC_DIRECTORY, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=%[1]s/%[2]s/peerOrganizations/%[2]s/peers/peer0.%[2]s/tls/server.key", constants.DEFAULT_FABRIC_DIRECTORY, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=%[1]s/%[2]s/peerOrganizations/%[2]s/peers/peer0.%[2]s/tls/ca.crt", constants.DEFAULT_FABRIC_DIRECTORY, domain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=%[1]s/%[2]s/peerOrganizations/%[2]s/users/Admin@%[2]s/msp", constants.DEFAULT_FABRIC_DIRECTORY, domain)),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode(constants.DEFAULT_FABRIC_DIRECTORY),
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
