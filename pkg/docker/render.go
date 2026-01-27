package docker

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type Renderer struct {
	config  pkg.Config
	path    string
	network string
}

func NewRenderer(config pkg.Config, network string, outputPath string) *Renderer {
	return &Renderer{
		config:  config,
		network: network,
		path:    outputPath,
	}
}

func RenderNetwork(networkName string, path string) (string, error) {

	if networkName == "" {
		networkName = "hyperledger_fabric_network"
	}

	err := yaml.MappingNode(
		yaml.ScalarNode("networks"),
		NewBridgeNetwork(networkName),
	).ToFile(fmt.Sprintf("%s/network.yml", path))

	return networkName, err
}

func (r *Renderer) RenderOrderers() error {
	var _orderers []*yaml.Node
	for _, orderer := range r.config.Orderers {
		node := NewOrderer(fmt.Sprintf("%s.%s", orderer.Hostname, orderer.Domain)).
			WithPort(orderer.Port).
			WithNetworks([]*yaml.Node{yaml.ScalarNode(r.network)})
		for _, n := range node.Content {
			_orderers = append(_orderers, (*yaml.Node)(n))
		}

	}

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		yaml.MappingNode(_orderers...),
	).ToFile(fmt.Sprintf("%s/orderer.yml", r.path))
}

func (r *Renderer) RenderPeerBase() error {
	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		NewPeerBase(r.network).
			Build(),
	).ToFile(fmt.Sprintf("%s/orgs/peer.base.yml", r.path))
}

func (r *Renderer) RenderCertificateAuthority(organization pkg.Organization) error {
	var nodes []*yaml.Node

	node := NewCertificateAuthority(fmt.Sprintf("ca.%s", organization.Domain)).
		WithNetworks([]*yaml.Node{yaml.ScalarNode(r.network)})

	if organization.CertificateAuthority.ExposePort > 0 {
		node.WithPort(organization.CertificateAuthority.ExposePort)
	}

	for _, n := range node.Content {
		nodes = append(nodes, (*yaml.Node)(n))
	}

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		yaml.MappingNode(nodes...),
	).ToFile(fmt.Sprintf("%s/orgs/%s/ca.yml", r.path, organization.Domain))
}

func (r *Renderer) RenderPeer(organization pkg.Organization, corePeerGossipBootstrap string, index int) error {
	node := NewPeer(
		fmt.Sprintf("%sMSP", organization.Name),
		fmt.Sprintf("peer%d.%s", index, organization.Domain),
		organization.Domain,
		corePeerGossipBootstrap,
		r.network,
	).Build()

	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		node,
	).ToFile(fmt.Sprintf("%s/orgs/%s/peer%d.yml", r.path, organization.Domain, index))
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

		if err := r.RenderPeers(organization); err != nil {
			return err
		}

		if err := r.RenderTools(organization); err != nil {
			return fmt.Errorf("Error when rendering the tools for the organization %s: %w", organization.Name, err)
		}
	}

	return nil
}

func (r *Renderer) RenderTools(organization pkg.Organization) error {
	return yaml.MappingNode(
		yaml.ScalarNode("services"),
		NewTools(
			organization.Name,
			organization.Domain,
			fmt.Sprintf("peer0.%s:7051", organization.Domain),
			fmt.Sprintf("%sMSP", organization.Name),
			r.network).
			Build(),
	).ToFile(fmt.Sprintf("%s/orgs/%s/tools.yml", r.path, organization.Domain))
}

func Render(config pkg.Config, path string) error {

	network, err := RenderNetwork(config.Docker.NetworkName, path)

	if err != nil {
		return fmt.Errorf("Error when rendering the network: %w", err)
	}

	renderer := NewRenderer(config, network, path)

	if err := renderer.RenderOrderers(); err != nil {
		return fmt.Errorf("Error when rendering the orderers: %w", err)
	}

	if err := renderer.RenderPeerBase(); err != nil {
		return fmt.Errorf("Error when rendering the peer base: %w", err)
	}

	if err := renderer.RenderOrganizations(); err != nil {
		return err
	}

	return nil
}
