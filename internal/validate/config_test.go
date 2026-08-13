package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"testing"
)

func TestConfigReturnsAllValidationErrors(t *testing.T) {
	configuration := spec.Config{
		Output: "invalid:name",
		Capabilities: spec.Capabilities{
			Channel:     "V1_4",
			Application: "V1_4",
			Orderer:     "V1_4",
		},
	}

	rules := make(map[RuleID]bool)
	for _, validationError := range Errors(Config(configuration)) {
		rules[validationError.RuleID] = true
	}
	want := []RuleID{
		RuleOutputDirectoryNameInvalid,
		RuleOrganizationsRequired,
		RuleChannelCapabilityUnsupported,
		RuleApplicationCapabilityUnsupported,
		RuleOrdererCapabilityUnsupported,
		RuleOrdererTopologyRequired,
	}

	for _, ruleID := range want {
		if !rules[ruleID] {
			t.Errorf("missing validation error for rule %s", ruleID)
		}
	}
}

func TestNewRulesAreCollected(t *testing.T) {
	c := spec.Config{Network: "bad name", Organizations: []spec.Organization{{Name: "Org1", Domain: "bad_domain", Users: -1, Peers: []spec.Peer{{Name: "p", Subdomain: "p"}, {Name: "p", Subdomain: "p"}}, Orderers: []spec.Orderer{{Name: "o", Subdomain: "o"}, {Name: "o", Subdomain: "o"}}}, {Name: "Org2", Domain: "bad_domain"}}, Profiles: []spec.Profile{{Name: "P"}, {Name: "P"}}, Channels: []spec.Channel{{Name: "channel", Profile: "Missing", Chaincodes: []spec.Chaincode{{Name: "cc"}, {Name: "cc"}}}, {Name: "channel"}}}
	want := []RuleID{RuleNetworkNameInvalid, RuleDomainInvalid, RuleOrganizationUsersInvalid, RuleOrganizationDomainDuplicate, RulePeerNameDuplicate, RulePeerSubdomainDuplicate, RuleOrdererNameDuplicate, RuleProfileNameDuplicate, RuleChannelProfileUndefined, RuleChannelNameDuplicate, RuleChaincodeNameDuplicate}
	found := map[RuleID]bool{}
	for _, err := range Errors(Config(c)) {
		found[err.RuleID] = true
	}
	for _, id := range want {
		if !found[id] {
			t.Errorf("missing rule %s", id)
		}
	}
}
