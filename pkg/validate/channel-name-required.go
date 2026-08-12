package validate

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"

func EmptyChannelNameFn(channel spec.Channel) error {
	if channel.Name == "" {
		return &ValidationError{
			RuleID: RuleChannelNameRequired,
			Rule:   "Empty Channel Name",
			Detail: "channel name cannot be empty",
		}
	}
	return nil
}
