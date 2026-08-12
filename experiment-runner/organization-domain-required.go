package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var organizationDomainRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleOrganizationDomainRequired,
	Apply:  func(node *yaml.Node) { organization(node, "Org1").GetValue("domain").SetScalar("", yaml.StringType) },
}}
