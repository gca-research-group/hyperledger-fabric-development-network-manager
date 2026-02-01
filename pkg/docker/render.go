package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type Renderer struct {
	config pkg.Config
}

func NewRenderer(config pkg.Config) *Renderer {
	if config.Network == "" {
		config.Network = "hyperledger_fabric_network"
	}

	return &Renderer{
		config: config,
	}
}

func (r *Renderer) RenderNetwork(networkName string, path string) error {
	return yaml.MappingNode(
		yaml.ScalarNode("networks"),
		NewBridgeNetwork(networkName),
	).ToFile(fmt.Sprintf("%s/network.yml", path))
}

func (r *Renderer) RenderOrderers(organization pkg.Organization) error {
	for _, orderer := range organization.Orderers {
		node := NewOrderer(fmt.Sprintf("%s.%s", orderer.Hostname, organization.Domain), organization.Domain).
			WithPort(orderer.Port).
			WithNetworks([]*yaml.Node{yaml.ScalarNode(r.config.Network)})

		err := yaml.MappingNode(
			yaml.ScalarNode("services"),
			node.Build(),
		).ToFile(fmt.Sprintf("%s/%s/%s.yml", r.config.Output, organization.Domain, orderer.Hostname))

		return err
	}

	return nil
}

func (r *Renderer) RenderPeerBase() error {
	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		NewPeerBase(r.config.Network).
			Build(),
	).ToFile(fmt.Sprintf("%s/peer.base.yml", r.config.Output))
}

func (r *Renderer) RenderCertificateAuthority(organization pkg.Organization) error {
	var nodes []*yaml.Node

	node := NewCertificateAuthority(fmt.Sprintf("ca.%s", organization.Domain)).
		WithNetworks([]*yaml.Node{yaml.ScalarNode(r.config.Network)})

	if organization.CertificateAuthority.ExposePort > 0 {
		node.WithPort(organization.CertificateAuthority.ExposePort)
	}

	for _, n := range node.Content {
		nodes = append(nodes, (*yaml.Node)(n))
	}

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		yaml.MappingNode(nodes...),
	).ToFile(fmt.Sprintf("%s/%s/ca.yml", r.config.Output, organization.Domain))
}

func (r *Renderer) RenderPeer(organization pkg.Organization, corePeerGossipBootstrap string, index int) error {
	node := NewPeer(
		fmt.Sprintf("%sMSP", organization.Name),
		fmt.Sprintf("peer%d.%s", index, organization.Domain),
		organization.Domain,
		corePeerGossipBootstrap,
		r.config.Network,
	).Build()

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		node,
	).ToFile(fmt.Sprintf("%s/%s/peer%d.yml", r.config.Output, organization.Domain, index))
}

func (r *Renderer) RenderPeers(organization pkg.Organization) error {
	peers := organization.Peers

	if peers == 0 {
		peers = 1
	}

	for i := 0; i < peers; i++ {
		corePeerGossipBootstrap := fmt.Sprintf("peer0.%s:7051", organization.Domain)

		if peers != 1 && i == 0 {
			corePeerGossipBootstrap = fmt.Sprintf("peer1.%s:7051", organization.Domain)
		} else {
			corePeerGossipBootstrap = fmt.Sprintf("peer0.%s:7051", organization.Domain)
		}

		if err := r.RenderPeer(organization, corePeerGossipBootstrap, i); err != nil {
			return fmt.Errorf("Error when rendering the peer %d for the organization %s: %w", i, organization.Name, err)
		}
	}

	return nil
}

func (r *Renderer) RenderOrganizations() error {
	for _, organization := range r.config.Organizations {
		if err := r.RenderCertificateAuthority(organization); err != nil {
			return fmt.Errorf("Error when rendering the certificate authority for the organization %s: %w", organization.Name, err)
		}

		if len(organization.Orderers) > 0 {
			if err := r.RenderOrderers(organization); err != nil {
				return fmt.Errorf("Error when rendering the orderers: %w", err)
			}
		}

		if err := r.RenderPeers(organization); err != nil {
			return err
		}

		var domains []string

		for _, o := range r.config.Organizations {
			if o.Domain != organization.Domain {
				domains = append(domains, o.Domain)
			}
		}

		if err := r.RenderTools(organization, domains); err != nil {
			return fmt.Errorf("Error when rendering the cryptomaterial.yml file for the organization %s: %w", organization.Name, err)
		}
	}

	return nil
}

func (r *Renderer) RenderTools(organization pkg.Organization, domains []string) error {
	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		NewTools(
			organization.Name,
			organization.Domain,
			fmt.Sprintf("peer0.%s:7051", organization.Domain),
			fmt.Sprintf("%sMSP", organization.Name),
			r.config.Network).Build(),
	).ToFile(fmt.Sprintf("%s/%s/tools.yml", r.config.Output, organization.Domain))
}

func (r *Renderer) RenderToolsWithMSP(currentOrganization pkg.Organization) error {
	var peerDomains []string
	var organizations []pkg.Organization

	for _, organization := range r.config.Organizations {
		if organization.Domain != currentOrganization.Domain {
			peerDomains = append(peerDomains, organization.Domain)
			if len(organization.Orderers) > 0 {
				organizations = append(organizations, organization)
			}
		}
	}

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		NewTools(
			currentOrganization.Name,
			currentOrganization.Domain,
			fmt.Sprintf("peer0.%s:7051", currentOrganization.Domain),
			fmt.Sprintf("%sMSP", currentOrganization.Name),
			r.config.Network).
			WithPeerMSPs(peerDomains).
			WithOrdererMSPs(organizations).
			Build(),
	).ToFile(fmt.Sprintf("%s/%s/tools.yml", r.config.Output, currentOrganization.Domain))
}

func (r *Renderer) Render() error {

	if err := r.RenderNetwork(r.config.Network, r.config.Output); err != nil {
		return fmt.Errorf("Error when rendering the network: %w", err)
	}

	if err := r.RenderPeerBase(); err != nil {
		return fmt.Errorf("Error when rendering the peer base: %w", err)
	}

	if err := r.RenderOrganizations(); err != nil {
		return err
	}

	return nil
}
