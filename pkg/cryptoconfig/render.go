package cryptoconfig

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"
)

func Render(config pkg.Config, path string) error {

	orderers := []*yaml.Node{}

	for _, orderer := range config.Orderers {
		orderers = append(orderers, NewOrdererOrg(orderer).Build())
	}

	organizations := []*yaml.Node{}

	for _, organization := range config.Organizations {
		organizations = append(organizations, NewPeerOrg(organization).Build())
	}

	return yaml.MappingNode(
		yaml.ScalarNode("OrdererOrgs"),
		yaml.SequenceNode(orderers...),
		yaml.ScalarNode("PeerOrgs"),
		yaml.SequenceNode(organizations...),
	).ToFile(fmt.Sprintf("%s/crypto-config.yml", path))
}
