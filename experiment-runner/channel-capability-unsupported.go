package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

func unsupportedChannelCapability(value string) MutationOperator {
	return MutationOperator{RuleID: validate.RuleChannelCapabilityUnsupported, Apply: func(node *yaml.Node) {
		node.GetValue("capabilities").GetValue("channel").SetScalar(value, yaml.StringType)
	}}
}

var channelCapabilityUnsupportedOperators = []MutationOperator{
	unsupportedChannelCapability(""), unsupportedChannelCapability("."), unsupportedChannelCapability("abc"),
}
