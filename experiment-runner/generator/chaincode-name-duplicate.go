package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var chaincodeNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleChaincodeNameDuplicate, Apply: func(n *yaml.Node) {
	chaincode(n, "defaultchannel", "Product").GetValue("name").SetScalar("Asset", yaml.StringType)
}}}
