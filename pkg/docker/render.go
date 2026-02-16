package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/constants"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type Renderer struct {
	config *pkg.Config
}

func NewRenderer(config *pkg.Config) *Renderer {
	if config.Network == "" {
		config.Network = constants.DEFAULT_NETORK
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
		node := NewOrderer(orderer.Subdomain, organization.Domain, r.config.Organizations).
			WithNetworks([]*yaml.Node{yaml.ScalarNode(r.config.Network)})

		err := yaml.MappingNode(
			yaml.ScalarNode("services"),
			node.Build(),
		).ToFile(fmt.Sprintf("%s/%s/%s.yml", r.config.Output, organization.Domain, orderer.Subdomain))

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

	node := NewCertificateAuthority(organization).
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

func (r *Renderer) RenderPeer(organization pkg.Organization, corePeerGossipBootstrap string, peer pkg.Peer) error {
	node := NewPeer(
		fmt.Sprintf("%sMSP", organization.Name),
		fmt.Sprintf("%s.%s", peer.Subdomain, organization.Domain),
		organization.Domain,
		corePeerGossipBootstrap,
		r.config.Network,
		r.config.Organizations,
	).Build()

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		node,
	).ToFile(fmt.Sprintf("%s/%s/%s.yml", r.config.Output, organization.Domain, peer.Subdomain))
}

func (r *Renderer) RenderPeers(organization pkg.Organization) error {

	for i, peer := range organization.Peers {
		gossipPeerIndex := 0

		if len(organization.Peers) != 1 && i == 0 {
			gossipPeerIndex = 1
		}

		gossipPeer := organization.Peers[gossipPeerIndex]

		gossipPeerport := gossipPeer.Port

		if gossipPeerport == 0 {
			gossipPeerport = constants.DEFAULT_PEER_PORT
		}

		corePeerGossipBootstrap := fmt.Sprintf("%s.%s:%d", gossipPeer.Subdomain, organization.Domain, gossipPeerport)

		if err := r.RenderPeer(organization, corePeerGossipBootstrap, peer); err != nil {
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
			organization,
			r.config.Organizations,
			r.config.Network).Build(),
	).ToFile(fmt.Sprintf("%s/%s/tools.yml", r.config.Output, organization.Domain))
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
