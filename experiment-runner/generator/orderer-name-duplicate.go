package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleOrdererNameDuplicate, Apply: func(n *yaml.Node) {
	orderers := organization(n, "Org1").GetValue("orderers")
	orderers.Content = append(orderers.Content, orderers.Content[0])
}}}
