package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

type Renderer struct {
	config *pkg.Config
}

func NewRenderer(config *pkg.Config) *Renderer {
	return &Renderer{config}
}

func (r *Renderer) Render() error {
	cfg, err := NewBuilder(r.config).Build()

	if err != nil {
		return err
	}

	return cfg.ToFile(fmt.Sprintf("%s/configtx.yml", r.config.Output))
}
