package validate

import (
	"errors"
	"testing"
)

func assertValidationError(t *testing.T, err error, ruleID RuleID, rule, detail string) {
	t.Helper()
	var validationError *ValidationError
	if !errors.As(err, &validationError) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if validationError.RuleID != ruleID || validationError.Rule != rule || validationError.Detail != detail {
		t.Fatalf("unexpected validation error: %#v", validationError)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
