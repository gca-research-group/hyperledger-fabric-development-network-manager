package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

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
