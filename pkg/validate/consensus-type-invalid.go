package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidConsensusTypeFn(profile spec.Profile) error {
	consensusType := profile.Consensus.Type

	if consensusType != "" && consensusType != "etcdraft" && consensusType != "BFT" {
		return &ValidationError{
			RuleID: RuleConsensusTypeInvalid,
			Rule:   "Invalid Consensus Type",
			Detail: fmt.Sprintf("invalid consensus type for the profile %s", profile.Name),
		}
	}

	return nil
}
