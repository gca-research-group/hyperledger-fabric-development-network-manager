package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var peerSubdomainRequiredOperators = []MutationOperator{{
	RuleID: validate.RulePeerSubdomainRequired,
	Apply: func(node *yaml.Node) {
		peer(node, "Org1", "Peer0").GetValue("subdomain").SetScalar("", yaml.StringType)
	},
}}
