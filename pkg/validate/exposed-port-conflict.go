package validate

import "fmt"

func ExposedPortConflictFn(port int, owner portOwner, exposedPorts map[int]portOwner) error {
	if port <= 0 {
		return nil
	}

	if existingOwner, exists := exposedPorts[port]; exists {
		return &ValidationError{
			RuleID: RuleExposedPortConflict,
			Rule:   "Exposed Port Conflict",
			Detail: fmt.Sprintf("Port %d is assigned to both %s and %s.", port, existingOwner, owner),
		}
	}

	exposedPorts[port] = owner

	return nil
}
