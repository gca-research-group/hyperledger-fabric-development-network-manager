package configtx

import "github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"

func capability(version string) *yaml.Node {
	node := yaml.MappingNode(
		yaml.ScalarNode(version),
		yaml.ScalarNode("true"),
	)

	return node
}

func NewApplicationCapability(version string) (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(ApplicationKey), capability(version).WithAnchor(ApplicationCapabilitiesKey)
}

func NewOrdererCapability(version string) (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(OrdererKey), capability(version).WithAnchor(OrdererCapabilitiesKey)
}

func NewChannelCapability(version string) (*yaml.Node, *yaml.Node) {
	return yaml.ScalarNode(ChannelKey), capability(version).WithAnchor(ChannelCapabilitiesKey)
}
