package fabric

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/directory"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

type Fabric struct {
	config  config.Config
	network string

	executor Executor

	identityManager *IdentityManager
}

func (f *Fabric) CleanUp() error {
	fmt.Print(">> Cleaning output folder...\n")
	return directory.RemoveFolderIfExists(f.config.Output)
}

func NewFabric(config config.Config, executor Executor) (*Fabric, error) {
	if len(config.Organizations) == 0 {
		return nil, fmt.Errorf("configuration must contain at least one organization")
	}

	if executor == nil {
		executor = &DefaultExecutor{}
	}

	network := fmt.Sprintf("%s/network.yml", config.Output)

	identityManager := NewIdentityManager(config, executor)

	return &Fabric{
		config,
		network,
		executor,
		identityManager,
	}, nil
}

func (f *Fabric) Start() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Start Certificate Authorities", f.RunCAContainers},
		{"Start Orderers", f.RunOrdererContainers},
		{"Start Peers", f.RunPeerContainers},
	}

	for _, step := range steps {
		fmt.Printf(">>> Step: %s\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}

func (f *Fabric) DeployNetwork() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Start Certificate Authorities", f.RunCAContainers},
		{"Generate Certificates", f.identityManager.GenerateAll},

		{"Generate Genesis", f.GenerateGenesisBlock},

		{"Start Orderers", f.RunOrdererContainers},
		{"Start Peers", f.RunPeerContainers},
		{"Join Orderers", f.JoinOrdererToTheChannel},
		{"Fetch Genesis Block", f.FetchGenesisBlock},
		{"Join Peers", f.JoinPeersToTheChannels},
	}

	for _, step := range steps {
		fmt.Printf(">>> Step: %s\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}
