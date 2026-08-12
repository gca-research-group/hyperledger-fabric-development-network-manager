package config

import "github.com/gca-research-group/fabric-network-orchestrator/pkg/validate"

// ValidateConfig is retained for compatibility. New code should use validate.Config.
func ValidateConfig(configuration Config) error {
	return validate.Config(configuration)
}
