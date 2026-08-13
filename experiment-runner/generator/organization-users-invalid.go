package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var organizationUsersInvalidOperators = []MutationOperator{{RuleID: validate.RuleOrganizationUsersInvalid, Apply: func(n *yaml.Node) {
	organization(n, "Org1").GetOrCreateValue("users", yaml.ScalarNode("0")).SetScalar("-1", yaml.IntType)
}}}
