package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var networkNameRequiredOperators = []MutationOperator{{RuleID: validate.RuleNetworkNameRequired, Apply: func(n *yaml.Node) { n.GetValue("network").SetScalar("", yaml.StringType) }}}
