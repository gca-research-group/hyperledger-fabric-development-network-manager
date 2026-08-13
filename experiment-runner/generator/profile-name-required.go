package generator

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var profileNameRequiredOperators = []MutationOperator{{RuleID: validate.RuleProfileNameRequired, Apply: func(n *yaml.Node) { profile(n, "DefaultProfile").GetValue("name").SetScalar("", yaml.StringType) }}}
