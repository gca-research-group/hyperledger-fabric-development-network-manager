package validate

import (
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func TestInvalidOutputDirectoryNameFn(t *testing.T) {
	assertValidationError(t, InvalidOutputDirectoryNameFn(spec.Config{Output: "foo:bar"}), RuleOutputDirectoryNameInvalid, "Invalid Output Directory Name", "the output directory must have a valid directory name")
	assertNoError(t, InvalidOutputDirectoryNameFn(spec.Config{Output: "./generated/network"}))
}

func TestIsValidOutputDirectory(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{path: "output", want: true},
		{path: "./output", want: true},
		{path: "generated/network", want: true},
		{path: `generated\network`, want: true},
		{path: "", want: false},
		{path: ".", want: false},
		{path: "../output", want: false},
		{path: "foo:bar", want: false},
		{path: "generated/bad?name", want: false},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			if got := IsValidOutputDirectory(test.path); got != test.want {
				t.Fatalf("IsValidOutputDirectory(%q) = %v, want %v", test.path, got, test.want)
			}
		})
	}
}
