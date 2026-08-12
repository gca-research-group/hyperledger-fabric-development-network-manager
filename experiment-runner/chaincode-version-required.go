package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var chaincodeVersionRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleChaincodeVersionRequired,
	Apply: func(node *yaml.Node) {
		chaincode(node, "defaultchannel", "Asset").GetValue("version").SetScalar("", yaml.StringType)
	},
}}
