package validate

import (
	"fmt"
	"regexp"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

var channelNameRegex = regexp.MustCompile(`^[a-z][a-z0-9.-]{0,248}$`)

func InvalidChannelNameFn(channel spec.Channel) error {
	if !channelNameRegex.MatchString(channel.Name) {
		return &ValidationError{
			RuleID: RuleChannelNameInvalid,
			Rule:   "Invalid Channel Name",
			Detail: fmt.Sprintf("invalid channel name: %s (must be lowercase alphanumeric, start with a letter, and contain only '.', '-', or alphanumeric characters, max 249 characters)", channel.Name),
		}
	}
	return nil
}
