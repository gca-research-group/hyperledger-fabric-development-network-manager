package cryptoconfig

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

type Renderer struct {
	config pkg.Config
}

func NewRenderer(config pkg.Config) *Renderer {
	return &Renderer{config}
}

func (r *Renderer) Render() error {

	for _, organization := range r.config.Organizations {
		orderers := []*yaml.Node{}
		peerOrgnization := NewPeerOrg(organization).Build()

		for _, orderer := range organization.Orderers {
			orderers = append(orderers, NewOrdererOrg(organization.Domain, orderer).Build())
		}

		var err error

		outputPath := fmt.Sprintf("%s/%s/crypto-config.yml", r.config.Output, organization.Domain)

		if len(orderers) > 0 {
			err = yaml.MappingNode(
				yaml.ScalarNode("OrdererOrgs"),
				yaml.SequenceNode(orderers...),
				yaml.ScalarNode("PeerOrgs"),
				yaml.SequenceNode([]*yaml.Node{peerOrgnization}...),
			).ToFile(outputPath)
		} else {
			err = yaml.MappingNode(
				yaml.ScalarNode("PeerOrgs"),
				yaml.SequenceNode([]*yaml.Node{peerOrgnization}...),
			).ToFile(outputPath)
		}

		if err != nil {
			return fmt.Errorf("Error when rendering the crypto-config.yml for the organization %s: %w", organization.Name, err)
		}
	}

	return nil
}
