package generator

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"

type IncompatibilityPolicy map[validate.RuleID][]validate.RuleID

func GenerateCombinations(groups [][]MutationOperator, combinationSize int, policy IncompatibilityPolicy) [][]MutationOperator {
	var result [][]MutationOperator

	if combinationSize <= 0 || len(groups) == 0 {
		return result
	}

	nonEmptyGroups := make([][]MutationOperator, 0, len(groups))

	for _, group := range groups {
		if len(group) > 0 {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	if len(nonEmptyGroups) == 0 {
		return result
	}

	effectiveSize := combinationSize

	if effectiveSize > len(nonEmptyGroups) {
		effectiveSize = len(nonEmptyGroups)
	}

	current := make([]MutationOperator, 0, effectiveSize)

	var combine func(startGroup int)

	combine = func(startGroup int) {

		if len(current) == effectiveSize {
			combination := make([]MutationOperator, len(current))
			copy(combination, current)

			result = append(result, combination)
			return
		}

		remaining := effectiveSize - len(current)

		if len(nonEmptyGroups)-startGroup < remaining {
			return
		}

		for i := startGroup; i < len(nonEmptyGroups); i++ {
			for _, element := range nonEmptyGroups[i] {
				if conflictsWithCurrent(element, current, policy) {
					continue
				}
				current = append(current, element)

				combine(i + 1)

				// Backtracking.
				current = current[:len(current)-1]
			}
		}
	}

	combine(0)

	return result
}

func conflictsWithCurrent(candidate MutationOperator, current []MutationOperator, policy IncompatibilityPolicy) bool {
	for _, selected := range current {
		if rulesConflict(candidate.RuleID, selected.RuleID, policy) {
			return true
		}
	}
	return false
}

func rulesConflict(first, second validate.RuleID, policy IncompatibilityPolicy) bool {
	return containsRule(policy[first], second) || containsRule(policy[second], first)
}

func containsRule(rules []validate.RuleID, target validate.RuleID) bool {
	for _, rule := range rules {
		if rule == target {
			return true
		}
	}
	return false
}
