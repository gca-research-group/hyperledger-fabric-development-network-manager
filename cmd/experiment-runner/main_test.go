package main

import (
	"strings"
	"testing"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/experiment/generator"
)

func TestParseOptions(t *testing.T) {
	options, err := parseOptions([]string{"--seed", "seed.yml"})
	if err != nil {
		t.Fatal(err)
	}
	if options.seedYAMLPath != "seed.yml" {
		t.Fatalf("seed path = %q, expected seed.yml", options.seedYAMLPath)
	}
	if options.mutationCount != generator.DefaultMutationCount {
		t.Fatalf("mutation count = %d, expected %d", options.mutationCount, generator.DefaultMutationCount)
	}
	if options.outputDirectory != "output" {
		t.Fatalf("output directory = %q, expected output", options.outputDirectory)
	}
}

func TestParseOptionsOverridesDefaults(t *testing.T) {
	options, err := parseOptions([]string{
		"--seed", "custom.yml",
		"--mutation-count", "2",
		"--output", "custom-output",
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.mutationCount != 2 || options.outputDirectory != "custom-output" {
		t.Fatalf("unexpected options: %+v", options)
	}
}

func TestParseOptionsRequiresSeedYAML(t *testing.T) {
	_, err := parseOptions(nil)
	if err == nil || !strings.Contains(err.Error(), "--seed") {
		t.Fatalf("expected missing seed YAML error, got %v", err)
	}
}

func TestParseOptionsRejectsInvalidMutationCount(t *testing.T) {
	_, err := parseOptions([]string{"--seed", "seed.yml", "--mutation-count", "invalid"})
	if err == nil || !strings.Contains(err.Error(), "invalid value") {
		t.Fatalf("expected invalid mutation count error, got %v", err)
	}
}

func TestRunReportsUnreadableSeedYAML(t *testing.T) {
	err := run([]string{"--seed", "missing-seed.yml"}, &strings.Builder{})
	if err == nil || !strings.Contains(err.Error(), "read seed YAML") {
		t.Fatalf("expected seed read error, got %v", err)
	}
}

func TestReportProgressThrottlesAndReportsFinalCount(t *testing.T) {
	var output strings.Builder
	for completed := 0; completed <= 2501; completed++ {
		reportProgress(&output, "generation", completed, 2501)
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
	reportProgress(&output, "validation", 0, 0)
	if output.String() != "validation progress: 0/0 (100.0%)\n" {
		t.Fatalf("unexpected zero-total progress: %q", output.String())
	}
}
