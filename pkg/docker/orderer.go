package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type OrdererNode struct {
	name string
	*yaml.Node
}

func NewOrderer(hostname string, currentOrganization config.Organization, organizations []config.Organization) *OrdererNode {
	ordererDomain := fmt.Sprintf("%s.%s", hostname, currentOrganization.Domain)

	ordererHostDir := fmt.Sprintf("./%[1]s/certificates/organizations/ordererOrganizations/%[1]s/orderers/%[2]s", currentOrganization.Domain, ordererDomain)
	ordererContainerDir := "/var/hyperledger/orderer"

	volumes := []*yaml.Node{
		yaml.ScalarNode(fmt.Sprintf("%s/msp:%s/msp", ordererHostDir, ordererContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("%s/tls:%s/tls", ordererHostDir, ordererContainerDir)),
	}

	cas := "/var/hyperledger/orderer/tls/ca.crt"

	version := currentOrganization.Version.Orderer

	if version == "" {
		version = constants.DEFAULT_FABRIC_VERSION
	}

	node := yaml.MappingNode(
		yaml.ScalarNode(ordererDomain),
		yaml.MappingNode(
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-orderer:%s", version)),
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(ordererDomain),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/var/hyperledger/orderer"),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode("ORDERER_GENERAL_LOGLEVEL=INFO"),
				yaml.ScalarNode("ORDERER_GENERAL_LISTENADDRESS=0.0.0.0"),
				yaml.ScalarNode("ORDERER_GENERAL_BOOTSTRAPMETHOD=none"),
				yaml.ScalarNode("ORDERER_GENERAL_LOCALMSPID=OrdererMSP"),
				yaml.ScalarNode("ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_ENABLED=true"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key"),
				yaml.ScalarNode("ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt"),
				yaml.ScalarNode(fmt.Sprintf("ORDERER_GENERAL_TLS_ROOTCAS=[%s]", cas)),
				yaml.ScalarNode("ORDERER_ADMIN_LISTENADDRESS=0.0.0.0:7053"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_ENABLED=true"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key"),
				yaml.ScalarNode("ORDERER_ADMIN_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt"),
				yaml.ScalarNode(fmt.Sprintf("ORDERER_ADMIN_TLS_CLIENTROOTCAS=[%s]", cas)),
				yaml.ScalarNode("ORDERER_CHANNELPARTICIPATION_ENABLED=true"),
			),
			yaml.ScalarNode("volumes"),
			yaml.SequenceNode(volumes...),
		),
	)

	return &OrdererNode{ordererDomain, node}
}

func (o *OrdererNode) WithNetworks(nodes []*yaml.Node) *OrdererNode {
	node := o.GetValue(o.name)
	node.GetOrCreateValue("networks", yaml.SequenceNode(nodes...))
	return o
}

// TODO:
/* func (o *OrdererNode) WithPort(port int) *OrdererNode {
	node := o.GetValue(o.name)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:7050", port))))
	return o
} */

func (o *OrdererNode) Build() *yaml.Node {
	return o.Node
}
