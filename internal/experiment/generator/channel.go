package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"strings"
)

var channelNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleChannelNameDuplicate, Apply: func(n *yaml.Node) {
	channels := n.GetValue("channels")
	channels.Content = append(channels.Content, channels.Content[0])
}}}

func invalidChannelName(value string) MutationOperator {
	return MutationOperator{RuleID: validate.RuleChannelNameInvalid, Apply: func(node *yaml.Node) {
		channel(node, "defaultchannel").GetValue("name").SetScalar(value, yaml.StringType)
	}}
}

var channelNameInvalidOperators = []MutationOperator{
	invalidChannelName("DefaultChannel"),
	invalidChannelName("1channel"),
	invalidChannelName("default channel"),
	invalidChannelName("a" + strings.Repeat("b", 249)),
}

var channelNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleChannelNameRequired,
	Apply:  func(node *yaml.Node) { channel(node, "defaultchannel").GetValue("name").SetScalar("", yaml.StringType) },
}}

var channelProfileRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleChannelProfileRequired,
	Apply: func(node *yaml.Node) {
		channel(node, "defaultchannel").GetValue("profile").SetScalar("", yaml.StringType)
	},
}}

var channelProfileUndefinedOperators = []MutationOperator{{RuleID: validate.RuleChannelProfileUndefined, Apply: func(n *yaml.Node) {
	channel(n, "defaultchannel").GetValue("profile").SetScalar("MissingProfile", yaml.StringType)
}}}
