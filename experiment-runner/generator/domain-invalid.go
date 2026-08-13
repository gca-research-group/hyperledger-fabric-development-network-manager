package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var domainInvalidOperators = []MutationOperator{{RuleID: validate.RuleDomainInvalid, Apply: func(n *yaml.Node) {
	organization(n, "Org1").GetValue("domain").SetScalar("invalid_domain", yaml.StringType)
}}}
