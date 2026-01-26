package configtx

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/yaml"

func defaultCapability() *yaml.Node {
	node := yaml.MappingNode(
		yaml.ScalarNode("V2_0"),
		yaml.ScalarNode("true"),
	)

	return node
}

func NewApplicationCapability() (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(ApplicationKey), defaultCapability().WithAnchor(ApplicationCapabilitiesKey)
}

func NewOrdererCapability() (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(OrdererKey), defaultCapability().WithAnchor(OrdererCapabilitiesKey)
}

func NewChannelCapability() (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(ChannelKey), defaultCapability().WithAnchor(ChannelCapabilitiesKey)
}
