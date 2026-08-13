package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var networkNameInvalidOperators = []MutationOperator{{RuleID: validate.RuleNetworkNameInvalid, Apply: func(n *yaml.Node) { n.GetValue("network").SetScalar("bad network", yaml.StringType) }}}
