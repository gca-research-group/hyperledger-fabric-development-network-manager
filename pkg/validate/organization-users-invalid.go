package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidOrganizationUsersFn(o spec.Organization) error {
	if o.Users < 0 {
		return validationError(RuleOrganizationUsersInvalid, "Invalid Organization Users", fmt.Sprintf("users of organization %s cannot be negative", o.Name))
	}
	return nil
}
