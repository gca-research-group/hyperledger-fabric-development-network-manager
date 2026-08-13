package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicateOrganizationDomainFn(o spec.Organization, seen map[string]struct{}) error {
	return duplicateValue(o.Domain, seen, RuleOrganizationDomainDuplicate, "Duplicate Organization Domain", fmt.Sprintf("duplicate organization domain: %s", o.Domain))
}
