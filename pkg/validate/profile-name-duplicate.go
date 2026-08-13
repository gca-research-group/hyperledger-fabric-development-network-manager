package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicateProfileNameFn(p spec.Profile, seen map[string]struct{}) error {
	return duplicateValue(p.Name, seen, RuleProfileNameDuplicate, "Duplicate Profile Name", fmt.Sprintf("duplicate profile name: %s", p.Name))
}
