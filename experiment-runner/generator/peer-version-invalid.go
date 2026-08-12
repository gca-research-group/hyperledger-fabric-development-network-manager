package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var peerVersionInvalidOperators = []MutationOperator{{
	RuleID: validate.RulePeerVersionInvalid,
	Apply: func(node *yaml.Node) {
		peer(node, "Org1", "Peer0").GetOrCreateValue("version", yaml.ScalarNode("")).SetScalar("1.4.0", yaml.StringType)
	},
}}
