package generator

import (
	"strings"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

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
