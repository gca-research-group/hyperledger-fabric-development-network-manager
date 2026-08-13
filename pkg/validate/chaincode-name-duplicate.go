package validate

import (
	"fmt"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func DuplicateChaincodeNameFn(c spec.Chaincode, channel string, seen map[string]struct{}) error {
	return duplicateValue(c.Name, seen, RuleChaincodeNameDuplicate, "Duplicate Chaincode Name", fmt.Sprintf("duplicate chaincode name in channel %s: %s", channel, c.Name))
}
