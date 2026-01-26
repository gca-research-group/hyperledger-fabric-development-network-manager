package configtx

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

func Render(config pkg.Config, path string) error {
	cfg, err := NewBuilder(config).Build()

	if err != nil {
		return err
	}

	return cfg.ToFile(path)
}
