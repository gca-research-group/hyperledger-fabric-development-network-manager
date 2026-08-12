package validate

func MultipleBootstrapOrganizationsFn(count int) error {
	if count > 1 {
		return &ValidationError{
			RuleID: RuleBootstrapOrganizationsMultiple,
			Rule:   "Multiple Bootstrap Orgs",
			Detail: "exactly one bootstrap organization must be defined",
		}
	}

	return nil
}
