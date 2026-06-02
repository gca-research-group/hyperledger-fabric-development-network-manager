package config

import "fmt"

type ValidationError struct {
	Layer  string
	RuleID string
	Rule   string
	Detail string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("[FATAL] Configuration validation failed!\n[LAYER] %s\n[RULE]  %s - %s\nDetail: %s",
		e.Layer, e.RuleID, e.Rule, e.Detail)
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
