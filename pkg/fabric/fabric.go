package fabric

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/command"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/configtx"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/cryptoconfig"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/docker"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/internal/directory"
)

type Fabric struct {
	config  pkg.Config
	network string

	crytpoConfigRenderer *cryptoconfig.Renderer
	configTxRenderer     *configtx.Renderer
	dockerRenderer       *docker.Renderer

	executor command.Executor

	identityManager *IdentityManager
}

func (f *Fabric) CleanUp() error {
	fmt.Print(">> Cleaning output folder...\n")
	return directory.RemoveFolderIfExists(f.config.Output)
}

func (f *Fabric) RenderConfigFiles() error {
	if err := f.crytpoConfigRenderer.Render(); err != nil {
		return err
	}

	if err := f.configTxRenderer.Render(); err != nil {
		return err
	}

	if err := f.dockerRenderer.Render(); err != nil {
		return err
	}

	return nil
}

func NewFabric(config pkg.Config, executor command.Executor) (*Fabric, error) {
	if len(config.Organizations) == 0 {
		return nil, fmt.Errorf("configuration must contain at least one organization")
	}

	if executor == nil {
		executor = &command.DefaultExecutor{}
	}

	network := fmt.Sprintf("%s/network.yml", config.Output)

	crytpoConfigRenderer := cryptoconfig.NewRenderer(config)
	configTxRenderer := configtx.NewRenderer(config)
	dockerRenderer := docker.NewRenderer(config)

	identityManager := NewIdentityManager(config, executor)

	return &Fabric{
		config,
		network,
		crytpoConfigRenderer,
		configTxRenderer,
		dockerRenderer,
		executor,
		identityManager,
	}, nil
}

func (f *Fabric) DeployNetwork() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Check Docker", f.IsDockerRunning},
		{"Clean Workspace", f.CleanUp},
		{"Render Config Files", f.RenderConfigFiles},
		{"Remove Old Containers", f.RemoveContainers},

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
