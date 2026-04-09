package compose

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
)

type PeerNode struct {
	*yaml.Node
	peer   config.Peer
	domain string
}

func NewPeer(
	peer config.Peer,
	currentOrganization config.Organization,
	corePeerGossipBootstrap string,
	network string,
	organizations []config.Organization,
) *PeerNode {

	domain := currentOrganization.Domain
	peerDomain := ResolvePeerDomain(peer.Subdomain, domain)
	peerPort := ResolvePeerPort(peer.Port)

	node := yaml.MappingNode(
		yaml.ScalarNode(peerDomain),
		yaml.MappingNode(
			yaml.ScalarNode("container_name"),
			yaml.ScalarNode(peerDomain),
			yaml.ScalarNode("image"),
			yaml.ScalarNode(fmt.Sprintf("hyperledger/fabric-peer:%s", ResolvePeerVersion(peer.Version))),
			yaml.ScalarNode("environment"),
			yaml.SequenceNode(
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", config.ResolveOrganizationMSPID(currentOrganization))),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ID=%s", peerDomain)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_ADDRESS=%s:%d", peerDomain, peerPort)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_LISTENADDRESS=0.0.0.0:%d", peerPort)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s", corePeerGossipBootstrap)),
				yaml.ScalarNode(fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:%d", peerDomain, peerPort)),
				yaml.ScalarNode("CORE_PEER_GOSSIP_USELEADERELECTION=true"),
				yaml.ScalarNode("CORE_PEER_GOSSIP_ORGLEADER=false"),
				yaml.ScalarNode("CORE_PEER_GOSSIP_STATE_ENABLED=true"),
				yaml.ScalarNode("CORE_LEDGER_STATE_STATEDATABASE=CouchDB"),
				yaml.ScalarNode(fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb.%s:5984", peerDomain)),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin"),
				yaml.ScalarNode("CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw"),
				yaml.ScalarNode("CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock"),
				yaml.ScalarNode(fmt.Sprintf("CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=%s", network)),
				yaml.ScalarNode("FABRIC_LOGGING_SPEC=INFO"),
				yaml.ScalarNode("CORE_PEER_PROFILE_ENABLED=true"),
				yaml.ScalarNode("CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp"),
				yaml.ScalarNode("CORE_PEER_TLS_ENABLED=true"),
				yaml.ScalarNode("CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key"),
				yaml.ScalarNode("CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt"),
				yaml.ScalarNode("CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt"),
				yaml.ScalarNode("CORE_PEER_TLS_CLIENTAUTHREQUIRED=false"),
				yaml.ScalarNode("CORE_PEER_TLS_CLIENTCERT_FILE=/etc/hyperledger/fabric/tls/server.crt"),
				yaml.ScalarNode("CORE_PEER_TLS_CLIENTKEY_FILE=/etc/hyperledger/fabric/tls/server.key"),
			),
			yaml.ScalarNode("working_dir"),
			yaml.ScalarNode("/etc/hyperledger/fabric"),
			yaml.ScalarNode("command"),
			yaml.ScalarNode("peer node start"),
			yaml.ScalarNode("networks"),
			yaml.SequenceNode(yaml.ScalarNode(network)),
		),
	)

	return &PeerNode{node, peer, domain}
}

func (pn *PeerNode) ExposePort() *PeerNode {
	if pn.peer.ExposePort == 0 {
		return pn
	}

	domain := pn.domain
	peerDomain := ResolvePeerDomain(pn.peer.Subdomain, domain)

	node := pn.GetValue(peerDomain)
	node.GetOrCreateValue("ports", yaml.SequenceNode(yaml.ScalarNode(fmt.Sprintf("%d:%d", pn.peer.ExposePort, ResolvePeerPort(pn.peer.Port)))))

	return pn
}

func (pn *PeerNode) WithVolumes() *PeerNode {
	domain := pn.domain
	peerDomain := ResolvePeerDomain(pn.peer.Subdomain, domain)

	peerHostDir := fmt.Sprintf("./%[1]s/certificate-authority/organizations/peerOrganizations/%[1]s/peers/%[2]s", domain, peerDomain)
	peerContainerDir := "/etc/hyperledger/fabric"

	volumes := []*yaml.Node{
		yaml.ScalarNode("/var/run/docker.sock:/host/var/run/docker.sock"),
		yaml.ScalarNode(fmt.Sprintf("%s/msp:%s/msp", peerHostDir, peerContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("%s/tls:%s/tls", peerHostDir, peerContainerDir)),
		yaml.ScalarNode(fmt.Sprintf("./%s/peers/%s/peer:/var/hyperledger/production", domain, pn.peer.Subdomain)),
	}

	node := pn.GetValue(peerDomain)
	node.GetOrCreateValue("volumes", yaml.SequenceNode(volumes...))

	return pn
}

func (pn *PeerNode) Build() *yaml.Node {
	return pn.Node
}
