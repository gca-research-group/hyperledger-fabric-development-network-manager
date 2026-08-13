package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var peerInternalPortInvalidOperators = []MutationOperator{{RuleID: validate.RulePeerInternalPortInvalid, Apply: func(n *yaml.Node) {
	peer(n, "Org1", "Peer0").GetOrCreateValue("port", yaml.ScalarNode("0")).SetScalar("65536", yaml.IntType)
}}}
