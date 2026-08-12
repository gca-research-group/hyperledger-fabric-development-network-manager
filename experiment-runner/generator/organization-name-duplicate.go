package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var organizationNameDuplicateOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationNameDuplicate,
	Apply:  func(node *yaml.Node) { organization(node, "Org2").GetValue("name").SetScalar("Org1", yaml.StringType) },
}}
