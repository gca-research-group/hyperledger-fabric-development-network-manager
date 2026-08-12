package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var profileOrganizationUndefinedOperators = []MutationOperator{{
	RuleID: validate.RuleProfileOrganizationUndefined,
	Apply: func(node *yaml.Node) {
		organizations := profile(node, "DefaultProfile").GetValue("organizations")
		organizations.Content[2].Value = "MissingOrg"
	},
}}
