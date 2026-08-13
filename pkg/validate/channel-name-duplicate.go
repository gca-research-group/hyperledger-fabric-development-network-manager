package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicateChannelNameFn(c spec.Channel, seen map[string]struct{}) error {
	return duplicateValue(c.Name, seen, RuleChannelNameDuplicate, "Duplicate Channel Name", fmt.Sprintf("duplicate channel name: %s", c.Name))
}
