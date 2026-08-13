package generator

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"

type IncompatibilityPolicy map[validate.RuleID][]validate.RuleID

func WalkCombinations(groups [][]MutationOperator, combinationSize int, policy IncompatibilityPolicy, visit func([]MutationOperator) error) error {
	if combinationSize <= 0 || len(groups) == 0 {
		return nil
	}

	nonEmptyGroups := make([][]MutationOperator, 0, len(groups))
	for _, group := range groups {
		if len(group) > 0 {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	if len(nonEmptyGroups) == 0 {
		return nil
	}

	effectiveSize := min(combinationSize, len(nonEmptyGroups))
	current := make([]MutationOperator, 0, effectiveSize)

	var combine func(int) error
	combine = func(startGroup int) error {
		if len(current) == effectiveSize {
			combination := append([]MutationOperator(nil), current...)
			return visit(combination)
		}

		remaining := effectiveSize - len(current)
		if len(nonEmptyGroups)-startGroup < remaining {
			return nil
		}

		for i := startGroup; i < len(nonEmptyGroups); i++ {
			for _, operator := range nonEmptyGroups[i] {
				if conflictsWithCurrent(operator, current, policy) {
					continue
				}

				current = append(current, operator)

				if err := combine(i + 1); err != nil {
					return err
				}

				current = current[:len(current)-1]
			}
		}

		return nil
	}

	return combine(0)
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
