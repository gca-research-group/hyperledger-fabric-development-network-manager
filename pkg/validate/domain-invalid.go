package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"regexp"
)

var domainRegex = regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+)$`)

func InvalidDomainFn(o spec.Organization) error {
	if o.Domain != "" && !domainRegex.MatchString(o.Domain) {
		return validationError(RuleDomainInvalid, "Invalid Domain", fmt.Sprintf("invalid organization domain: %s", o.Domain))
	}
	return nil
}
