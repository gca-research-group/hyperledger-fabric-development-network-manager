package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var ordererNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleOrdererNameDuplicate, Apply: func(n *yaml.Node) {
	orderers := organization(n, "Org1").GetValue("orderers")
	orderers.Content = append(orderers.Content, orderers.Content[0])
}}}

var ordererNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererNameRequired,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetValue("name").SetScalar("", yaml.StringType)
	},
}}

var ordererSubdomainRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererSubdomainRequired,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetValue("subdomain").SetScalar("", yaml.StringType)
	},
}}

var ordererVersionInvalidOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererVersionInvalid,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetOrCreateValue("version", yaml.ScalarNode("")).SetScalar("1.4.0", yaml.StringType)
	},
}}
