package cmd

import (
	"errors"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

var configPath string

func LoadConfig() (*pkg.Config, error) {
	if configPath == "" {
		return nil, errors.New("missing required flag: --config")
	}

	return pkg.LoadConfigFromPath(configPath)
}
