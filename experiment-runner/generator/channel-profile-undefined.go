package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var channelProfileUndefinedOperators = []MutationOperator{{RuleID: validate.RuleChannelProfileUndefined, Apply: func(n *yaml.Node) {
	channel(n, "defaultchannel").GetValue("profile").SetScalar("MissingProfile", yaml.StringType)
}}}
