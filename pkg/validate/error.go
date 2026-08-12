package validate

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	RuleID RuleID
	Rule   string
	Detail string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("[FATAL] Configuration validation failed!\n[RULE]  %s - %s\nDetail: %s", e.RuleID, e.Rule, e.Detail)
}

func Errors(err error) []*ValidationError {
	if err == nil {
		return nil
	}

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		var result []*ValidationError
		for _, child := range joined.Unwrap() {
			result = append(result, Errors(child)...)
		}
		return result
	}

	var validationError *ValidationError
	if errors.As(err, &validationError) {
		return []*ValidationError{validationError}
	}
	return nil
}

type portOwner struct {
	ownerType string
	name      string
}

func (po portOwner) String() string {
	if po.ownerType == "Certificate Authority" {
		return fmt.Sprintf("Certificate Authority of org '%s'", po.name)
	}

	return fmt.Sprintf("%s '%s'", po.ownerType, po.name)
}
