package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
)

var certificateAuthorityPortInvalidOperators = []MutationOperator{
	{
		RuleID: validate.RuleCertificateAuthorityPortInvalid,
		Apply: func(node *yaml.Node) {
			organizations := node.GetValue("organizations")

			org := organizations.FindByValue("name", "Org1")
			certificateAuthority := org.GetOrCreateValue("certificateAuthority", yaml.MappingNode())
			exposePort := certificateAuthority.GetOrCreateValue("exposePort", yaml.ScalarNode("0"))
			exposePort.SetScalar("-1", yaml.IntType)
		},
	},
}

var exposedPortConflictOperators = []MutationOperator{{
	RuleID: validate.RuleExposedPortConflict,
	Apply: func(node *yaml.Node) {
		org1Port := peer(node, "Org1", "Peer0").GetOrCreateValue("exposePort", yaml.ScalarNode("0"))
		peer(node, "Org2", "Peer0").GetOrCreateValue("exposePort", yaml.ScalarNode("0")).SetScalar(org1Port.Value, yaml.IntType)
	},
}}

var ordererInternalPortInvalidOperators = []MutationOperator{{RuleID: validate.RuleOrdererInternalPortInvalid, Apply: func(n *yaml.Node) {
	orderer(n, "Org1", "Orderer").GetOrCreateValue("port", yaml.ScalarNode("0")).SetScalar("65536", yaml.IntType)
}}}

var ordererPortInvalidOperators = []MutationOperator{{
	RuleID: validate.RuleOrdererPortInvalid,
	Apply: func(node *yaml.Node) {
		orderer(node, "Org1", "Orderer").GetOrCreateValue("exposePort", yaml.ScalarNode("0")).SetScalar("-1", yaml.IntType)
	},
}}

var peerInternalPortInvalidOperators = []MutationOperator{{RuleID: validate.RulePeerInternalPortInvalid, Apply: func(n *yaml.Node) {
	peer(n, "Org1", "Peer0").GetOrCreateValue("port", yaml.ScalarNode("0")).SetScalar("65536", yaml.IntType)
}}}

var peerPortInvalidOperators = []MutationOperator{{
	RuleID: validate.RulePeerPortInvalid,
	Apply: func(node *yaml.Node) {
		peer(node, "Org1", "Peer0").GetOrCreateValue("exposePort", yaml.ScalarNode("0")).SetScalar("-1", yaml.IntType)
	},
}}
