package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func EmptyProfileNameFn(p spec.Profile, i int) error {
	if p.Name == "" {
		return validationError(RuleProfileNameRequired, "Empty Profile Name", fmt.Sprintf("name of the profile index %d is empty", i))
	}
	return nil
}
