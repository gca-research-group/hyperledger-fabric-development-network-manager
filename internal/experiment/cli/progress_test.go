package cli

import (
	"strings"
	"testing"
)

func TestReportProgressThrottlesAndReportsFinalCount(t *testing.T) {
	var output strings.Builder
	for completed := 0; completed <= 2501; completed++ {
		reportProgress(&output, "generation", completed, 2501, defaultProgressInterval)
	}

	expected := "generation progress: 0/2501 (0.0%)\n" +
		"generation progress: 1000/2501 (40.0%)\n" +
		"generation progress: 2000/2501 (80.0%)\n" +
		"generation progress: 2501/2501 (100.0%)\n"
	if output.String() != expected {
		t.Fatalf("progress output:\n%s\nexpected:\n%s", output.String(), expected)
	}
}

func TestReportProgressHandlesZeroTotal(t *testing.T) {
	var output strings.Builder
	reportProgress(&output, "validation", 0, 0, defaultProgressInterval)
	if output.String() != "validation progress: 0/0 (100.0%)\n" {
		t.Fatalf("unexpected zero-total progress: %q", output.String())
	}
}

func TestProgressIntervalRejectsInvalidValues(t *testing.T) {
	for _, command := range []string{"generate", "validate"} {
		for _, value := range []string{"0", "-1", "invalid"} {
			args := []string{command, "--progress-interval", value}
			if command == "generate" {
				args = append(args, "--seed", "missing.yaml")
			}
			err := run(args, &strings.Builder{})
			if err == nil || !strings.Contains(err.Error(), "progress-interval") {
				t.Fatalf("%v: expected interval error before accessing files, got %v", args, err)
			}
		}
	}
}
