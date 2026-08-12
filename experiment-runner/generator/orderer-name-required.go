package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererNameRequired,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetValue("name").SetScalar("", yaml.StringType)
	},
}}
