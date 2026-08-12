package validate

func NoOrdererTopologyFn(hasOrderer bool) error {
	if !hasOrderer {
		return &ValidationError{
			RuleID: RuleOrdererTopologyRequired,
			Rule:   "No Orderer Topology",
			Detail: "at least one orderer must be configured",
		}
	}

	return nil
}
