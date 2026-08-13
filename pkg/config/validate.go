package config

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"

func ValidateConfig(configuration Config) error {
	return validate.Config(configuration)
}
