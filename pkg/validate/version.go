package validate

import (
	"fmt"
	"strings"
)

func validateBinary(version, minimum string) error {
	if version == "" {
		return nil
	}

	currentParts, minimumParts := parseVersion(version), parseVersion(minimum)

	for i := 0; i < 3; i++ {
		if currentParts[i] > minimumParts[i] {
			return nil
		}

		if currentParts[i] < minimumParts[i] {
			return fmt.Errorf("version %s is lower than required %s", version, minimum)
		}
	}
	return nil
}

func parseVersion(version string) [3]int {
	var result [3]int
	parts := strings.Split(version, ".")

	for i := 0; i < len(parts) && i < 3; i++ {
		fmt.Sscanf(parts[i], "%d", &result[i])
	}

	return result
}
