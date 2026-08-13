package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func UndefinedChannelProfileFn(c spec.Channel, profiles map[string]struct{}) error {
	if c.Profile != "" {
		if _, ok := profiles[c.Profile]; !ok {
			return validationError(RuleChannelProfileUndefined, "Undefined Channel Profile", fmt.Sprintf("profile not defined: %s", c.Profile))
		}
	}
	return nil
}
