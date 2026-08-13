package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var profileNameDuplicateOperators = []MutationOperator{{RuleID: validate.RuleProfileNameDuplicate, Apply: func(n *yaml.Node) {
	profiles := n.GetValue("profiles")
	profiles.Content = append(profiles.Content, profiles.Content[0])
}}}
