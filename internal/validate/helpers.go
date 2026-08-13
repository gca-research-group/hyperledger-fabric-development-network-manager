package validate

func validOptionalTCPPort(port int) bool { return port == 0 || port >= 1 && port <= 65535 }

func validationError(id RuleID, rule, detail string) error {
	return &ValidationError{RuleID: id, Rule: rule, Detail: detail}
}

func duplicateValue(value string, seen map[string]struct{}, id RuleID, rule, detail string) error {
	if value == "" {
		return nil
	}
	if _, exists := seen[value]; exists {
		return validationError(id, rule, detail)
	}
	seen[value] = struct{}{}
	return nil
}
