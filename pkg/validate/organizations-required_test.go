package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"testing"
)

func TestNoOrganizationFn(t *testing.T) {
	assertValidationError(t, NoOrganizationFn(spec.Config{}), RuleOrganizationsRequired, "No Organization", "at least one organization must be defined")
	assertNoError(t, NoOrganizationFn(spec.Config{Organizations: []spec.Organization{{Name: "Org1"}}}))
}
