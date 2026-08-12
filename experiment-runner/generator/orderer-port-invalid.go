package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererPortInvalidOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererPortInvalid,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetOrCreateValue("exposePort", yaml.ScalarNode("0")).SetScalar("-1", yaml.IntType)
	},
}}
