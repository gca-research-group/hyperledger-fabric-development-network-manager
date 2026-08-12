package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var exposedPortConflictOperators = []MutationOperator{{
	RuleID: validate.RuleExposedPortConflict,
	Apply: func(node *yaml.Node) {
		peer(node, "Org2", "Peer0").GetValue("exposePort").SetScalar("7051", yaml.IntType)
	},
}}
