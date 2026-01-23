package configtx

import (
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"gopkg.in/yaml.v3"
)

func Render(config pkg.Config) (*yaml.Node, error) {
	cfg, err := NewBuilder(config).Build()
	if err != nil {
		return nil, err
	}
	return cfg.MarshalYAML()
}
