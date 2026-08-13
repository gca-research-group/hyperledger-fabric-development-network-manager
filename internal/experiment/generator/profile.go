package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var consensusTypeInvalidOperators = []MutationOperator{{
	RuleID: validate.RuleConsensusTypeInvalid,
	Apply: func(node *yaml.Node) {
		consensus := profile(node, "DefaultProfile").GetOrCreateValue("consensus", yaml.MappingNode())
		consensus.GetOrCreateValue("type", yaml.ScalarNode("")).SetScalar("raft", yaml.StringType)
	},
}}

var profileNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleProfileNameDuplicate, Apply: func(n *yaml.Node) {
	profiles := n.GetValue("profiles")
	profiles.Content = append(profiles.Content, profiles.Content[0])
}}}

var profileNameRequiredOperators = []MutationOperator{{RuleID: validate.RuleProfileNameRequired, Apply: func(n *yaml.Node) { profile(n, "DefaultProfile").GetValue("name").SetScalar("", yaml.StringType) }}}

var profileOrganizationUndefinedOperators = []MutationOperator{{
	RuleID: validate.RuleProfileOrganizationUndefined,
	Apply: func(node *yaml.Node) {
		organizations := profile(node, "DefaultProfile").GetValue("organizations")
		organizations.Content[2].Value = "MissingOrg"
	},
}}

var profileOrganizationsRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleProfileOrganizationsRequired,
	Apply:  func(node *yaml.Node) { profile(node, "DefaultProfile").GetValue("organizations").Content = nil },
}}
