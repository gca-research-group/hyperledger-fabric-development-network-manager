package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
	yamlv3 "gopkg.in/yaml.v3"
)

var peerSubdomainDuplicateOperators = []MutationOperator{{RuleID: validate.RulePeerSubdomainDuplicate, Apply: func(n *yaml.Node) {
	peers := organization(n, "Org1").GetValue("peers")
	duplicate := (*yaml.Node)(peers.Content[0]).Clone()
	duplicate.GetValue("name").SetScalar("Peer1", yaml.StringType)
	duplicate.GetValue("exposePort").SetScalar("0", yaml.IntType)
	peers.Content = append(peers.Content, (*yamlv3.Node)(duplicate))
}}}
