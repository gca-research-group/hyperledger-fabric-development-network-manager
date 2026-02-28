package network

import (
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
	network := compose.ResolveNetworkDockerComposeFile(config.Output)

	identityManager := NewIdentityManager(config, exec)

	return &Network{
		config,
		network,
		exec,
		identityManager,
	}, nil
}
