package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"regexp"
)

var networkNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

func InvalidNetworkNameFn(c spec.Config) error {
	if c.Network != "" && !networkNameRegex.MatchString(c.Network) {
		return validationError(RuleNetworkNameInvalid, "Invalid Network Name", "network name must contain only letters, numbers, '.', '_', or '-' and start with a letter or number")
	}
	return nil
}
