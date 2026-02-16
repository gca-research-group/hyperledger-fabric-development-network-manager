package configtx

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type Renderer struct {
	config *config.Config
}

func NewRenderer(config *config.Config) *Renderer {
	return &Renderer{config}
}

func (r *Renderer) Render() error {
	cfg, err := NewBuilder(r.config).Build()

	if err != nil {
		return err
	}

	return cfg.ToFile(fmt.Sprintf("%s/configtx.yml", r.config.Output))
}
