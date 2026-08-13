package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"regexp"
)

func DuplicateChannelNameFn(c spec.Channel, seen map[string]struct{}) error {
	return duplicateValue(c.Name, seen, RuleChannelNameDuplicate, "Duplicate Channel Name", fmt.Sprintf("duplicate channel name: %s", c.Name))
}

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

func UndefinedChannelProfileFn(c spec.Channel, profiles map[string]struct{}) error {
	if c.Profile != "" {
		if _, ok := profiles[c.Profile]; !ok {
			return validationError(RuleChannelProfileUndefined, "Undefined Channel Profile", fmt.Sprintf("profile not defined: %s", c.Profile))
		}
	}
	return nil
}
