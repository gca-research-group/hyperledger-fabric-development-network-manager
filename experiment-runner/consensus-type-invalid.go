package main

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
)

var consensusTypeInvalidOperators = []MutationOperator{{
	RuleID: validate.RuleConsensusTypeInvalid,
	Apply: func(node *yaml.Node) {
		consensus := profile(node, "DefaultProfile").GetOrCreateValue("consensus", yaml.MappingNode())
		consensus.GetOrCreateValue("type", yaml.ScalarNode("")).SetScalar("raft", yaml.StringType)
	},
}}
