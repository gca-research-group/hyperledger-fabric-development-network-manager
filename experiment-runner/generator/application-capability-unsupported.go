package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
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
