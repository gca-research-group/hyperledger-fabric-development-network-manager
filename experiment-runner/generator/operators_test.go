package generator

import (
	"errors"
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"
	yamlv3 "gopkg.in/yaml.v3"
)

var allRuleIDs = []validate.RuleID{
	validate.RuleOutputDirectoryNameInvalid,
	validate.RuleOrganizationsRequired,
	validate.RuleOrganizationNameRequired,
	validate.RuleOrganizationDomainRequired,
	validate.RuleCertificateAuthorityPortInvalid,
	validate.RulePeerNameRequired,
	validate.RulePeerSubdomainRequired,
	validate.RulePeerPortInvalid,
	validate.RuleOrdererNameRequired,
	validate.RuleOrdererSubdomainRequired,
	validate.RuleOrdererPortInvalid,
	validate.RuleChaincodeNameRequired,
	validate.RuleChaincodePathRequired,
	validate.RuleChaincodeVersionRequired,
	validate.RuleProfileOrganizationsRequired,
	validate.RuleChannelProfileRequired,
	validate.RuleChannelNameRequired,
	validate.RuleChannelNameInvalid,
	validate.RuleChannelCapabilityUnsupported,
	validate.RuleApplicationCapabilityUnsupported,
	validate.RuleOrdererCapabilityUnsupported,
	validate.RuleOrganizationNameDuplicate,
	validate.RulePeerVersionInvalid,
	validate.RuleOrdererVersionInvalid,
	validate.RuleOrdererTopologyRequired,
	validate.RuleBootstrapOrganizationsMultiple,
	validate.RuleConsensusTypeInvalid,
	validate.RuleProfileOrganizationUndefined,
	validate.RuleExposedPortConflict,
}

func seedNode(t *testing.T) *yaml.Node {
	t.Helper()
	node, err := yaml.FromBytes([]byte(seedYAML))
	if err != nil {
		t.Fatalf("parse seed: %v", err)
	}
	return node
}

func decodeConfig(t *testing.T, node *yaml.Node) spec.Config {
	t.Helper()
	var configuration spec.Config
	if err := (*yamlv3.Node)(node.Document()).Decode(&configuration); err != nil {
		t.Fatalf("decode mutation: %v", err)
	}
	return configuration
}

func TestSeedIsValid(t *testing.T) {
	configuration := decodeConfig(t, seedNode(t))
	if err := validate.Config(configuration); err != nil {
		t.Fatalf("seed must be valid (decoded output %q): %v", configuration.Output, err)
	}
}

func TestEveryOperatorTriggersDeclaredRule(t *testing.T) {
	for _, group := range operators {
		for index, operator := range group {
			operator := operator
			t.Run(string(operator.RuleID)+"/"+string(rune('0'+index)), func(t *testing.T) {
				node := seedNode(t)
				operator.Apply(node.Document())
				err := validate.Config(decodeConfig(t, node))
				var validationError *validate.ValidationError
				if !errors.As(err, &validationError) {
					t.Fatalf("expected ValidationError, got %v", err)
				}
				if validationError.RuleID != operator.RuleID {
					t.Fatalf("operator declares %s, triggered %s", operator.RuleID, validationError.RuleID)
				}
			})
		}
	}
}

func TestOperatorRegistryCoversEveryRuleOnce(t *testing.T) {
	registered := make(map[validate.RuleID]int)
	for _, group := range operators {
		if len(group) == 0 {
			t.Fatal("operator registry contains an empty group")
		}
		groupRule := group[0].RuleID
		for _, operator := range group {
			if operator.RuleID != groupRule {
				t.Fatalf("group %s also contains %s", groupRule, operator.RuleID)
			}
		}
		registered[groupRule]++
	}
	for _, ruleID := range allRuleIDs {
		if registered[ruleID] != 1 {
			t.Errorf("rule %s registered %d times", ruleID, registered[ruleID])
		}
		delete(registered, ruleID)
	}
	for ruleID := range registered {
		t.Errorf("unexpected registered rule %s", ruleID)
	}
}

func TestGenerateCombinationsUsesAtMostOneOperatorPerRule(t *testing.T) {
	for _, combination := range GenerateCombinations(operators, 3) {
		seen := make(map[validate.RuleID]bool)
		for _, operator := range combination {
			if seen[operator.RuleID] {
				t.Fatalf("combination contains rule %s more than once", operator.RuleID)
			}
			seen[operator.RuleID] = true
		}
	}
}
