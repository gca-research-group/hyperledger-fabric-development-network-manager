package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererVersionInvalidOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererVersionInvalid,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetOrCreateValue("version", yaml.ScalarNode("")).SetScalar("1.4.0", yaml.StringType)
	},
}}
