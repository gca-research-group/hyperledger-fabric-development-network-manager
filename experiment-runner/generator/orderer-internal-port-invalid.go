package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererInternalPortInvalidOperators = []MutationOperator{{RuleID: validate.RuleOrdererInternalPortInvalid, Apply: func(n *yaml.Node) {
	orderer(n, "Org1", "Orderer").GetOrCreateValue("port", yaml.ScalarNode("0")).SetScalar("65536", yaml.IntType)
}}}
