package compose

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"

func NewBridgeNetwork(name string) *yaml.Node {
	return yaml.MappingNode(
		yaml.ScalarNode(name),
		yaml.MappingNode(
			yaml.ScalarNode("name"),
			yaml.ScalarNode(name),
			yaml.ScalarNode("driver"),
			yaml.ScalarNode("bridge"),
		),
	)
}
