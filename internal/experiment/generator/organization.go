package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var bootstrapOrganizationsMultipleOperators = []MutationOperator{
	{
		RuleID: validate.RuleBootstrapOrganizationsMultiple,
		Apply: func(node *yaml.Node) {
			organizations := node.GetValue("organizations")

			org1 := organizations.FindByValue("name", "Org1")
			bootstrapOrg1 := org1.GetOrCreateValue("bootstrap", yaml.ScalarNode("false"))
			bootstrapOrg1.SetScalar("true", yaml.BoolType)

			org2 := organizations.FindByValue("name", "Org2")
			bootstrapOrg2 := org2.GetOrCreateValue("bootstrap", yaml.ScalarNode("false"))
			bootstrapOrg2.SetScalar("true", yaml.BoolType)
		},
	},
}

var domainInvalidOperators = []MutationOperator{{RuleID: validate.RuleDomainInvalid, Apply: func(n *yaml.Node) {
	organization(n, "Org1").GetValue("domain").SetScalar("invalid_domain", yaml.StringType)
}}}

var organizationDomainDuplicateOperators = []MutationOperator{{RuleID: validate.RuleOrganizationDomainDuplicate, Apply: func(n *yaml.Node) {
	org1Domain := organization(n, "Org1").GetValue("domain")
	organization(n, "Org2").GetValue("domain").SetScalar(org1Domain.Value, yaml.StringType)
}}}

var organizationDomainRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationDomainRequired,
	Apply:  func(node *yaml.Node) { organization(node, "Org1").GetValue("domain").SetScalar("", yaml.StringType) },
}}

var organizationNameDuplicateOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationNameDuplicate,
	Apply:  func(node *yaml.Node) { organization(node, "Org2").GetValue("name").SetScalar("Org1", yaml.StringType) },
}}

var organizationNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationNameRequired,
	Apply:  func(node *yaml.Node) { organization(node, "Org1").GetValue("name").SetScalar("", yaml.StringType) },
}}

var organizationUsersInvalidOperators = []MutationOperator{{RuleID: validate.RuleOrganizationUsersInvalid, Apply: func(n *yaml.Node) {
	organization(n, "Org1").GetOrCreateValue("users", yaml.ScalarNode("0")).SetScalar("-1", yaml.IntType)
}}}

var organizationsRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationsRequired,
	Apply:  func(node *yaml.Node) { node.GetValue("organizations").Content = nil },
}}
