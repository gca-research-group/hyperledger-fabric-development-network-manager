package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererTopologyRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererTopologyRequired,
	Apply:  func(node *yaml.Node) { organization(node, "Org1").GetValue("orderers").Content = nil },
}}
