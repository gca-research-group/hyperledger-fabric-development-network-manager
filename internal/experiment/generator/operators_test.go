package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/validate"
	"github.com/gca-research-group/fabric-network-orchestrator/internal/yaml"
	yamlv3 "gopkg.in/yaml.v3"
)

var allRuleIDs = []validate.RuleID{
	validate.RuleOutputDirectoryNameInvalid,
	validate.RuleNetworkNameRequired,
	validate.RuleNetworkNameInvalid,
	validate.RuleOrganizationsRequired,
	validate.RuleOrganizationNameRequired,
	validate.RuleOrganizationDomainRequired,
	validate.RuleOrganizationDomainDuplicate,
	validate.RuleOrganizationUsersInvalid,
	validate.RuleDomainInvalid,
	validate.RuleCertificateAuthorityPortInvalid,
	validate.RulePeerNameRequired,
	validate.RulePeerSubdomainRequired,
	validate.RulePeerPortInvalid,
	validate.RulePeerInternalPortInvalid,
	validate.RulePeerNameDuplicate,
	validate.RulePeerSubdomainDuplicate,
	validate.RuleOrdererNameRequired,
	validate.RuleOrdererSubdomainRequired,
	validate.RuleOrdererPortInvalid,
	validate.RuleOrdererInternalPortInvalid,
	validate.RuleOrdererNameDuplicate,
	validate.RuleChaincodeNameRequired,
	validate.RuleChaincodePathRequired,
	validate.RuleChaincodeVersionRequired,
	validate.RuleChaincodeNameDuplicate,
	validate.RuleProfileOrganizationsRequired,
	validate.RuleProfileNameRequired,
	validate.RuleProfileNameDuplicate,
	validate.RuleChannelProfileRequired,
	validate.RuleChannelNameRequired,
	validate.RuleChannelNameInvalid,
	validate.RuleChannelNameDuplicate,
	validate.RuleChannelProfileUndefined,
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
	err := WalkCombinations(operators, 3, incompatibilities, func(combination []MutationOperator) error {
		seen := make(map[validate.RuleID]bool)
		for _, operator := range combination {
			if seen[operator.RuleID] {
				return fmt.Errorf("combination contains rule %s more than once", operator.RuleID)
			}
			seen[operator.RuleID] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWalkCombinationsStreamsOneHundredThousandCombinations(t *testing.T) {
	groups := make([][]MutationOperator, 5)
	for groupIndex := range groups {
		for operatorIndex := 0; operatorIndex < 10; operatorIndex++ {
			groups[groupIndex] = append(groups[groupIndex], MutationOperator{RuleID: validate.RuleID(fmt.Sprintf("rule-%d-%d", groupIndex, operatorIndex))})
		}
	}

	count := 0
	if err := WalkCombinations(groups, 5, nil, func([]MutationOperator) error {
		count++
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if count != 100_000 {
		t.Fatalf("visited %d combinations, expected 100000", count)
	}
}

func TestGenerateStreamsValidManifest(t *testing.T) {
	outputDirectory := t.TempDir()
	summary, err := Generate(outputDirectory, 1)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(outputDirectory, "scenarios.json"))
	if err != nil {
		t.Fatal(err)
	}
	var scenarios []ScenarioRules
	if err := json.Unmarshal(data, &scenarios); err != nil {
		t.Fatalf("decode generated manifest: %v", err)
	}
	if summary.Total != len(scenarios) || summary.Total == 0 {
		t.Fatalf("summary total %d does not match manifest length %d", summary.Total, len(scenarios))
	}
	for _, scenario := range scenarios {
		if _, err := os.Stat(filepath.Join(outputDirectory, "config", scenario.Scenario+".yaml")); err != nil {
			t.Fatalf("scenario %s was not written: %v", scenario.Scenario, err)
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
				count := 0
				if err := WalkCombinations(groups, 2, incompatibilities, func([]MutationOperator) error { count++; return nil }); err != nil {
					t.Fatal(err)
				}
				if count != 0 {
					t.Fatalf("generated rules %s and %s in an incompatible order", first, second)
				}
			}
		}
	}
}

func TestGeneratedCombinationsExcludeIncompatibleRules(t *testing.T) {
	represented := make(map[validate.RuleID]bool)
	err := WalkCombinations(operators, DefaultMutationCount, incompatibilities, func(combination []MutationOperator) error {
		for i, first := range combination {
			represented[first.RuleID] = true
			for _, second := range combination[i+1:] {
				if rulesConflict(first.RuleID, second.RuleID, incompatibilities) {
					return fmt.Errorf("generated incompatible rules %s and %s", first.RuleID, second.RuleID)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, ruleID := range allRuleIDs {
		if !represented[ruleID] {
			t.Errorf("rule %s is not represented by any generated combination", ruleID)
		}
	}
}

func TestGeneratedCombinationsTriggerEveryDeclaredRule(t *testing.T) {
	index := 0
	err := WalkCombinations(operators, DefaultMutationCount, incompatibilities, func(combination []MutationOperator) error {
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
				return fmt.Errorf("combination %d (%v) did not trigger declared rule %s", index, combinationRuleIDs(combination), operator.RuleID)
			}
		}
		index++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func combinationRuleIDs(combination []MutationOperator) []validate.RuleID {
	rules := make([]validate.RuleID, 0, len(combination))
	for _, operator := range combination {
		rules = append(rules, operator.RuleID)
	}
	return rules
}

func TestOrganizationsRequiredExcludesOrganizationDomainRequired(t *testing.T) {
	err := WalkCombinations(operators, DefaultMutationCount, incompatibilities, func(combination []MutationOperator) error {
		hasOrganizationsRequired := false
		hasOrganizationDomainRequired := false
		for _, operator := range combination {
			hasOrganizationsRequired = hasOrganizationsRequired || operator.RuleID == validate.RuleOrganizationsRequired
			hasOrganizationDomainRequired = hasOrganizationDomainRequired || operator.RuleID == validate.RuleOrganizationDomainRequired
		}
		if hasOrganizationsRequired && hasOrganizationDomainRequired {
			return errors.New("generated organizations.required with organization.domain.required")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
