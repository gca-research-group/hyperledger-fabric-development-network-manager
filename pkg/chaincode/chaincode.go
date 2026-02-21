package chaincode

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type Chaincode struct {
	config   *config.Config
	network  string
	executor executor.Executor
}

func NewChaincode(config *config.Config, exec executor.Executor) *Chaincode {
	if exec == nil {
		exec = &executor.DefaultExecutor{}
	}

	network := compose.ResolveNetworkDockerComposeFile(config.Output)

	return &Chaincode{config, network, exec}
}

func (c *Chaincode) Publish() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Package the Chaincodes", c.Package},
		{"Install in the Peers", c.Install},
		{"Approving the Chaincodes", c.Approve},
		{"Commit the Chaincodes", c.Commit},
	}

	for _, step := range steps {
		fmt.Printf(">>> Step: %s\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}
