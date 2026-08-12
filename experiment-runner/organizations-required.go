package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var organizationsRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationsRequired,
	Apply:  func(node *yaml.Node) { node.GetValue("organizations").Content = nil },
}}
