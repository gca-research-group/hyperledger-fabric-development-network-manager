package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var peerPortInvalidOperators = []MutationOperator{{
	RuleID: validate.RulePeerPortInvalid,
	Apply: func(node *yaml.Node) {
		peer(node, "Org1", "Peer0").GetValue("exposePort").SetScalar("-1", yaml.IntType)
	},
}}
