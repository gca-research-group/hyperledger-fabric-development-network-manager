package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var organizationDomainDuplicateOperators = []MutationOperator{{RuleID: validate.RuleOrganizationDomainDuplicate, Apply: func(n *yaml.Node) {
	organization(n, "Org2").GetValue("domain").SetScalar("org1.example.com", yaml.StringType)
}}}
