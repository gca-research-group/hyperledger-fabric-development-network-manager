package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var ordererTopologyRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererTopologyRequired,
	Apply:  func(node *yaml.Node) { organization(node, "Org1").GetValue("orderers").Content = nil },
}}
