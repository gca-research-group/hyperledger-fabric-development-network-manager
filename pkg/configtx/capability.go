package configtx

import "github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/yaml"

func capability(version string) *yaml.Node {
	node := yaml.MappingNode(
		yaml.ScalarNode(version),
		yaml.ScalarNode("true"),
	)

	return node
}

func NewApplicationCapability() (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(ApplicationKey), capability("V2_5").WithAnchor(ApplicationCapabilitiesKey)
}

func NewOrdererCapability() (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(OrdererKey), capability("V2_0").WithAnchor(OrdererCapabilitiesKey)
}

func NewChannelCapability() (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(ChannelKey), capability("V2_0").WithAnchor(ChannelCapabilitiesKey)
}
