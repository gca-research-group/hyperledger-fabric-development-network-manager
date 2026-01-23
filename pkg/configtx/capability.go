package configtx

func defaultCapability() *Node {
	node := MappingNode(
		ScalarNode("V2_0"),
		ScalarNode("true"),
	)

	return node
}

func NewApplicationCapability() (*Node, *Node) {
	return ScalarNode(ApplicationKey), defaultCapability().WithAnchor(ApplicationCapabilitiesKey)
}

func NewOrdererCapability() (*Node, *Node) {
	return ScalarNode(OrdererKey), defaultCapability().WithAnchor(OrdererCapabilitiesKey)
}

func NewChannelCapability() (*Node, *Node) {
	return ScalarNode(ChannelKey), defaultCapability().WithAnchor(ChannelCapabilitiesKey)
}
