package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	yamlv3 "gopkg.in/yaml.v3"
)

var peerNameDuplicateOperators = []MutationOperator{{RuleID: validate.RulePeerNameDuplicate, Apply: func(n *yaml.Node) {
	peers := organization(n, "Org1").GetValue("peers")
	duplicate := (*yaml.Node)(peers.Content[0]).Clone()
	duplicate.GetValue("subdomain").SetScalar("peer1", yaml.StringType)
	duplicate.GetOrCreateValue("exposePort", yaml.ScalarNode("0")).SetScalar("0", yaml.IntType)
	peers.Content = append(peers.Content, (*yamlv3.Node)(duplicate))
}}}

var peerNameRequiredOperators = []MutationOperator{{
	RuleID: validate.RulePeerNameRequired,
	Apply:  func(node *yaml.Node) { peer(node, "Org1", "Peer0").GetValue("name").SetScalar("", yaml.StringType) },
}}

var peerSubdomainDuplicateOperators = []MutationOperator{{RuleID: validate.RulePeerSubdomainDuplicate, Apply: func(n *yaml.Node) {
	peers := organization(n, "Org1").GetValue("peers")
	duplicate := (*yaml.Node)(peers.Content[0]).Clone()
	duplicate.GetValue("name").SetScalar("Peer1", yaml.StringType)
	duplicate.GetOrCreateValue("exposePort", yaml.ScalarNode("0")).SetScalar("0", yaml.IntType)
	peers.Content = append(peers.Content, (*yamlv3.Node)(duplicate))
}}}

var peerSubdomainRequiredOperators = []MutationOperator{{
	RuleID: validate.RulePeerSubdomainRequired,
	Apply: func(node *yaml.Node) {
		peer(node, "Org1", "Peer0").GetValue("subdomain").SetScalar("", yaml.StringType)
	},
}}

var peerVersionInvalidOperators = []MutationOperator{{
	RuleID: validate.RulePeerVersionInvalid,
	Apply: func(node *yaml.Node) {
		peer(node, "Org1", "Peer0").GetOrCreateValue("version", yaml.ScalarNode("")).SetScalar("1.4.0", yaml.StringType)
	},
}}
