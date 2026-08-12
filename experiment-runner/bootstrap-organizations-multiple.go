package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var bootstrapOrganizationsMultipleOperators = []MutationOperator{
	{
		RuleID: validate.RuleBootstrapOrganizationsMultiple,
		Apply: func(node *yaml.Node) {
			organizations := node.GetValue("organizations")

			org1 := organizations.FindByValue("name", "Org1")
			bootstrapOrg1 := org1.GetValue("bootstrap")
			bootstrapOrg1.SetScalar("true", yaml.BoolType)

			org2 := organizations.FindByValue("name", "Org2")
			bootstrapOrg2 := org2.GetValue("bootstrap")
			bootstrapOrg2.SetScalar("true", yaml.BoolType)
		},
	},
}
