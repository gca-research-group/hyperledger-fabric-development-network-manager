package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var networkNameInvalidOperators = []MutationOperator{{RuleID: validate.RuleNetworkNameInvalid, Apply: func(n *yaml.Node) { n.GetValue("network").SetScalar("bad network", yaml.StringType) }}}

var networkNameRequiredOperators = []MutationOperator{{RuleID: validate.RuleNetworkNameRequired, Apply: func(n *yaml.Node) { n.GetValue("network").SetScalar("", yaml.StringType) }}}

var outputDirectoryOperators = []MutationOperator{
	{
		RuleID: validate.RuleOutputDirectoryNameInvalid,
		Apply: func(node *yaml.Node) {
			output := node.GetValue("output")
			output.SetScalar("", yaml.StringType)
		},
	},
	{
		RuleID: validate.RuleOutputDirectoryNameInvalid,
		Apply: func(node *yaml.Node) {
			output := node.GetValue("output")
			output.SetScalar(".", yaml.StringType)
		},
	},
	{
		RuleID: validate.RuleOutputDirectoryNameInvalid,
		Apply: func(node *yaml.Node) {
			output := node.GetValue("output")
			output.SetScalar("..", yaml.StringType)
		},
	},
	{
		RuleID: validate.RuleOutputDirectoryNameInvalid,
		Apply: func(node *yaml.Node) {
			output := node.GetValue("output")
			output.SetScalar("foo:bar", yaml.StringType)
		},
	},
	{
		RuleID: validate.RuleOutputDirectoryNameInvalid,
		Apply: func(node *yaml.Node) {
			output := node.GetValue("output")
			output.SetScalar("generated/bad?name", yaml.StringType)
		},
	},
}
