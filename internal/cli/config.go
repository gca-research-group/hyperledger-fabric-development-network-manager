package cli

import (
	"errors"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

var configPath string

func LoadConfig() (*config.Config, error) {
	if configPath == "" {
		return nil, errors.New("missing required flag: --config")
	}

	return config.LoadConfigFromPath(configPath)
}
