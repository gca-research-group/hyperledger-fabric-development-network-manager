package network

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type Network struct {
	config  config.Config
	network string

	executor executor.Executor

	identityManager *IdentityManager
}

func NewNetwork(config config.Config, exec executor.Executor) (*Network, error) {
	if len(config.Organizations) == 0 {
		return nil, fmt.Errorf("configuration must contain at least one organization")
	}

	network := compose.ResolveNetworkDockerComposeFile(config.Output)

	identityManager := NewIdentityManager(config, exec)

	return &Network{
		config,
		network,
		exec,
		identityManager,
	}, nil
}
