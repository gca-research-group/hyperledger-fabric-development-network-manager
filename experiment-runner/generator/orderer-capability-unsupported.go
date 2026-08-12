package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

func unsupportedOrdererCapability(value string) MutationOperator {
	return MutationOperator{RuleID: validate.RuleOrdererCapabilityUnsupported, Apply: func(node *yaml.Node) {
		node.GetValue("capabilities").GetValue("orderer").SetScalar(value, yaml.StringType)
	}}
}

var ordererCapabilityUnsupportedOperators = []MutationOperator{
	unsupportedOrdererCapability(""), unsupportedOrdererCapability("."), unsupportedOrdererCapability("abc"),
}
