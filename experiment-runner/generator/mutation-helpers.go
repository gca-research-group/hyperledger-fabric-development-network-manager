package generator

import "github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"

func organization(node *yaml.Node, name string) *yaml.Node {
	return node.GetValue("organizations").FindByValue("name", name)
}

func peer(node *yaml.Node, organizationName, name string) *yaml.Node {
	return organization(node, organizationName).GetValue("peers").FindByValue("name", name)
}

func orderer(node *yaml.Node, organizationName, name string) *yaml.Node {
	return organization(node, organizationName).GetValue("orderers").FindByValue("name", name)
}

func profile(node *yaml.Node, name string) *yaml.Node {
	return node.GetValue("profiles").FindByValue("name", name)
}

func channel(node *yaml.Node, name string) *yaml.Node {
	channels := node.GetValue("channels")
	result := channels.FindByValue("name", name)
	if result == nil && len(channels.Content) > 0 {
		return (*yaml.Node)(channels.Content[0])
	}
	return result
}

func chaincode(node *yaml.Node, channelName, name string) *yaml.Node {
	chaincodes := channel(node, channelName).GetValue("chaincodes")
	result := chaincodes.FindByValue("name", name)
	if result == nil && len(chaincodes.Content) > 0 {
		return (*yaml.Node)(chaincodes.Content[0])
	}
	return result
}
