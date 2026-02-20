package compose

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type OrdererNode struct {
	*yaml.Node
	orderer config.Orderer
	domain  string
}

func NewOrderer(orderer config.Orderer, currentOrganization config.Organization, organizations []config.Organization) *OrdererNode {
	domain := currentOrganization.Domain
	ordererDomain := ResolveOrdererDomain(orderer.Subdomain, domain)

	cas := "/var/hyperledger/orderer/tls/ca.crt"

	version := ResolveOrdererVersion(currentOrganization.Version.Orderer)

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
		),
	)

	return &OrdererNode{node, orderer, domain}
}

func (o *OrdererNode) WithNetworks(nodes []*yaml.Node) *OrdererNode {
	node := o.GetValue(ResolveOrdererDomain(o.orderer.Subdomain, o.domain))
	node.GetOrCreateValue("networks", yaml.SequenceNode(nodes...))
	return o
}

func (o *OrdererNode) WithVolumes() *OrdererNode {
	domain := o.domain
	ordererDomain := ResolveOrdererDomain(o.orderer.Subdomain, domain)

	ordererHostDir := fmt.Sprintf("./%[1]s/certificate-authority/organizations/ordererOrganizations/%[1]s/orderers/%[2]s", o.domain, ordererDomain)
	ordererContainerDir := "/var/hyperledger/orderer"

	volumes := []*yaml.Node{
		yaml.ScalarNode(fmt.Sprintf("%s/msp:%s/msp", ordererHostDir, ordererContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("%s/tls:%s/tls", ordererHostDir, ordererContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("./%s/orderers/%s/orderer:/var/hyperledger/production/orderer", domain, o.orderer.Subdomain)),
	}

	node := o.GetValue(ordererDomain)
	node.GetOrCreateValue("volumes", yaml.SequenceNode(volumes...))

	return o
}

func (o *OrdererNode) ExposePort() *OrdererNode {
	if o.orderer.ExposePort == 0 {
		return o
	}

	ordererDomain := ResolveOrdererDomain(o.orderer.Subdomain, o.domain)

	node := o.GetValue(ordererDomain)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:%d", o.orderer.ExposePort, ResolveOrdererPort(o.orderer.Port)))))

	return o
}

func (o *OrdererNode) Build() *yaml.Node {
	return o.Node
}
