package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var channelNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleChannelNameDuplicate, Apply: func(n *yaml.Node) {
	channels := n.GetValue("channels")
	channels.Content = append(channels.Content, channels.Content[0])
}}}
