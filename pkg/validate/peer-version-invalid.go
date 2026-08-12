package validate

import (
	"fmt"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func InvalidPeerVersionFn(peer spec.Peer, organizationName, channelCapability string) error {
	if err := validateBinary(peer.Version, spec.MinBinaryVersion[channelCapability]); err != nil {
		return &ValidationError{
			RuleID: RulePeerVersionInvalid,
			Rule:   "Invalid Peer Version",
			Detail: fmt.Sprintf("peer version of org %s invalid: %v", organizationName, err),
		}
	}

	return nil
}
