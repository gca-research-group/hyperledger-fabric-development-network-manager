package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyChannelProfileFn(channel spec.Channel) error {
	if channel.Profile == "" {
		return &ValidationError{
			RuleID: RuleChannelProfileRequired,
			Rule:   "Empty Channel Profile",
			Detail: fmt.Sprintf("channel %s must reference a profile", channel.Name),
		}
	}
	return nil
}
