package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var applicationCapabilityUnsupportedOperators = []MutationOperator{
	{
		RuleID: validate.RuleApplicationCapabilityUnsupported,
		Apply: func(node *yaml.Node) {
			capabilities := node.GetValue("capabilities")
			application := capabilities.GetValue("application")
			application.SetScalar("", yaml.StringType)
		},
	},
	{
		RuleID: validate.RuleApplicationCapabilityUnsupported,
		Apply: func(node *yaml.Node) {
			capabilities := node.GetValue("capabilities")
			application := capabilities.GetValue("application")
			application.SetScalar(".", yaml.StringType)
		},
	},
	{
		RuleID: validate.RuleApplicationCapabilityUnsupported,
		Apply: func(node *yaml.Node) {
			capabilities := node.GetValue("capabilities")
			application := capabilities.GetValue("application")
			application.SetScalar("abc", yaml.StringType)
		},
	},
}

func unsupportedChannelCapability(value string) MutationOperator {
	return MutationOperator{RuleID: validate.RuleChannelCapabilityUnsupported, Apply: func(node *yaml.Node) {
		node.GetValue("capabilities").GetValue("channel").SetScalar(value, yaml.StringType)
	}}
}

var channelCapabilityUnsupportedOperators = []MutationOperator{
	unsupportedChannelCapability(""), unsupportedChannelCapability("."), unsupportedChannelCapability("abc"),
}

func unsupportedOrdererCapability(value string) MutationOperator {
	return MutationOperator{RuleID: validate.RuleOrdererCapabilityUnsupported, Apply: func(node *yaml.Node) {
		node.GetValue("capabilities").GetValue("orderer").SetScalar(value, yaml.StringType)
	}}
}

var ordererCapabilityUnsupportedOperators = []MutationOperator{
	unsupportedOrdererCapability(""), unsupportedOrdererCapability("."), unsupportedOrdererCapability("abc"),
}
