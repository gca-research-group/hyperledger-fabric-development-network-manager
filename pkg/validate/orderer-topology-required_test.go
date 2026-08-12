package validate

import "testing"

func TestNoOrdererTopologyFn(t *testing.T) {
	assertValidationError(t, NoOrdererTopologyFn(false), RuleOrdererTopologyRequired, "No Orderer Topology", "at least one orderer must be configured")
	assertNoError(t, NoOrdererTopologyFn(true))
}
