package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var ordererSubdomainRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererSubdomainRequired,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetValue("subdomain").SetScalar("", yaml.StringType)
	},
}}
