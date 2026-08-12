package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var peerNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RulePeerNameRequired,
	Apply:  func(node *yaml.Node) { peer(node, "Org1", "Peer0").GetValue("name").SetScalar("", yaml.StringType) },
}}
