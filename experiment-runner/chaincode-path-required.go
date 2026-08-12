package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var chaincodePathRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleChaincodePathRequired,
	Apply: func(node *yaml.Node) {
		chaincode(node, "defaultchannel", "Asset").GetValue("path").SetScalar("", yaml.StringType)
	},
}}
