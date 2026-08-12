package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var certificateAuthorityPortInvalidOperators = []MutationOperator{
	{
		RuleID: validate.RuleCertificateAuthorityPortInvalid,
		Apply: func(node *yaml.Node) {
			organizations := node.GetValue("organizations")

			org := organizations.FindByValue("name", "Org1")
			certificateAuthority := org.GetValue("certificateAuthority")
			exposePort := certificateAuthority.GetValue("exposePort")
			exposePort.SetScalar("-1", yaml.IntType)
		},
	},
}
