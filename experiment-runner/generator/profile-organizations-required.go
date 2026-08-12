package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var profileOrganizationsRequiredOperators = []MutationOperator{{
	RuleID: validate.RuleProfileOrganizationsRequired,
	Apply:  func(node *yaml.Node) { profile(node, "DefaultProfile").GetValue("organizations").Content = nil },
}}
