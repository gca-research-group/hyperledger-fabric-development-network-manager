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
	for _, combination := range GenerateCombinations(operators, 3, incompatibilities) {
		seen := make(map[validate.RuleID]bool)
		for _, operator := range combination {
			if seen[operator.RuleID] {
				t.Fatalf("combination contains rule %s more than once", operator.RuleID)
			}
			seen[operator.RuleID] = true
		}
	}
}

func TestIncompatibilitiesAreSymmetric(t *testing.T) {
	for first, incompatibleRules := range incompatibilities {
		for _, second := range incompatibleRules {
			if !rulesConflict(first, second, incompatibilities) || !rulesConflict(second, first, incompatibilities) {
				t.Fatalf("rules %s and %s must conflict in both directions", first, second)
			}

			firstOperator := MutationOperator{RuleID: first}
			secondOperator := MutationOperator{RuleID: second}
			for _, groups := range [][][]MutationOperator{
				{{firstOperator}, {secondOperator}},
				{{secondOperator}, {firstOperator}},
			} {
				if combinations := GenerateCombinations(groups, 2, incompatibilities); len(combinations) != 0 {
					t.Fatalf("generated rules %s and %s in an incompatible order", first, second)
				}
			}
		}
	}
}

func TestGeneratedCombinationsExcludeIncompatibleRules(t *testing.T) {
	combinations := GenerateCombinations(operators, DefaultMutationCount, incompatibilities)
	represented := make(map[validate.RuleID]bool)

	for _, combination := range combinations {
		for i, first := range combination {
			represented[first.RuleID] = true
			for _, second := range combination[i+1:] {
				if rulesConflict(first.RuleID, second.RuleID, incompatibilities) {
					t.Fatalf("generated incompatible rules %s and %s", first.RuleID, second.RuleID)
				}
			}
		}
	}

	for _, ruleID := range allRuleIDs {
		if !represented[ruleID] {
			t.Errorf("rule %s is not represented by any generated combination", ruleID)
		}
	}
}

func TestGeneratedCombinationsTriggerEveryDeclaredRule(t *testing.T) {
	for index, combination := range GenerateCombinations(operators, DefaultMutationCount, incompatibilities) {
		node := seedNode(t)
		for _, operator := range combination {
			operator.Apply(node.Document())
		}

		actual := make(map[validate.RuleID]bool)
		for _, validationError := range validate.Errors(validate.Config(decodeConfig(t, node))) {
			actual[validationError.RuleID] = true
		}
		for _, operator := range combination {
			if !actual[operator.RuleID] {
				t.Fatalf("combination %d did not trigger declared rule %s", index, operator.RuleID)
			}
		}
	}
}

func TestOrganizationsRequiredExcludesOrganizationDomainRequired(t *testing.T) {
	for _, combination := range GenerateCombinations(operators, DefaultMutationCount, incompatibilities) {
		hasOrganizationsRequired := false
		hasOrganizationDomainRequired := false
		for _, operator := range combination {
			hasOrganizationsRequired = hasOrganizationsRequired || operator.RuleID == validate.RuleOrganizationsRequired
			hasOrganizationDomainRequired = hasOrganizationDomainRequired || operator.RuleID == validate.RuleOrganizationDomainRequired
		}
		if hasOrganizationsRequired && hasOrganizationDomainRequired {
			t.Fatal("generated organizations.required with organization.domain.required")
		}
	}
}
