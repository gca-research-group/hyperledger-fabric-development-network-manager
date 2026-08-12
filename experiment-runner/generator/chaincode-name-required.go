package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var chaincodeNameRequiredOperators = []MutationOperator{
	{
		RuleID: validate.RuleChaincodeNameRequired,
		Apply: func(node *yaml.Node) {
			chaincode(node, "defaultchannel", "Asset").GetValue("name").SetScalar("", yaml.StringType)
		},
	},
}
