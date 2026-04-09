package network

import (
	"github.com/gca-research-group/fabric-network-orchestrator/internal/executor"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/compose"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
)

type Network struct {
	config  config.Config
	network string

	executor executor.Executor

	identityManager  *IdentityManager
	containerManager *ContainerManager
}

func NewNetwork(config config.Config, exec executor.Executor) *Network {
	network := compose.ResolveNetworkDockerComposeFile(config.Output)

	identityManager := NewIdentityManager(config, exec)
	containerManager := NewContainerManager(config, exec)

	return &Network{
		config,
		network,
		exec,
		identityManager,
		containerManager,
	}
}
