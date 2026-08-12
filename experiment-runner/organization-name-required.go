package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var organizationNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationNameRequired,
	Apply:  func(node *yaml.Node) { organization(node, "Org1").GetValue("name").SetScalar("", yaml.StringType) },
}}
