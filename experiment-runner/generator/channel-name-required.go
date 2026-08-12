package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var channelNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleChannelNameRequired,
	Apply:  func(node *yaml.Node) { channel(node, "defaultchannel").GetValue("name").SetScalar("", yaml.StringType) },
}}
