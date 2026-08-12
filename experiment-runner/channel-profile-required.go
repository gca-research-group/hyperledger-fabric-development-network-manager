package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var channelProfileRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleChannelProfileRequired,
	Apply: func(node *yaml.Node) {
		channel(node, "defaultchannel").GetValue("profile").SetScalar("", yaml.StringType)
	},
}}
