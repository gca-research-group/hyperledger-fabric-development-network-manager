package validate

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/spec"
	"path/filepath"
	"regexp"
	"strings"
)

var networkNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

func InvalidNetworkNameFn(c spec.Config) error {
	if c.Network != "" && !networkNameRegex.MatchString(c.Network) {
		return validationError(RuleNetworkNameInvalid, "Invalid Network Name", "network name must contain only letters, numbers, '.', '_', or '-' and start with a letter or number")
	}
	return nil
}

func EmptyNetworkNameFn(c spec.Config) error {
	if c.Network == "" {
		return validationError(RuleNetworkNameRequired, "Empty Network Name", "network name cannot be empty")
	}
	return nil
}

var validDirName = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func IsValidDirectoryName(name string) bool {
	return name != "" && name != "." && name != ".." && validDirName.MatchString(name)
}

func IsValidOutputDirectory(path string) bool {
	if path == "" {
		return false
	}

	path = strings.TrimPrefix(path, filepath.VolumeName(path))
	parts := strings.FieldsFunc(path, func(r rune) bool { return r == '/' || r == '\\' })

	hasDirectoryName := false

	for _, part := range parts {
		if part == "." {
			continue
		}

		if !IsValidDirectoryName(part) {
			return false
		}

		hasDirectoryName = true
	}

	return hasDirectoryName
}

func InvalidOutputDirectoryNameFn(configuration spec.Config) error {
	if !IsValidOutputDirectory(configuration.Output) {
		return &ValidationError{
			RuleID: RuleOutputDirectoryNameInvalid,
			Rule:   "Invalid Output Directory Name",
			Detail: "the output directory must have a valid directory name",
		}
	}

	return nil
}
