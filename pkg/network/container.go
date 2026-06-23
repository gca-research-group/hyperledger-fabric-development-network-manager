package network

import (
	"fmt"
	"log/slog"

	"github.com/gca-research-group/fabric-network-orchestrator/internal/executor"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/compose"
	"github.com/gca-research-group/fabric-network-orchestrator/pkg/config"
)

type ContainerManager struct {
	config   config.Config
	executor executor.Executor
	network  string
}

func NewContainerManager(config config.Config, executor executor.Executor) *ContainerManager {
	network := compose.ResolveNetworkDockerComposeFile(config.Output)

	return &ContainerManager{
		config:   config,
		executor: executor,
		network:  network,
	}
}

func (cm *ContainerManager) RunOrdererContainers() error {
	slog.Info("Executing orderer containers")
	for _, organization := range cm.config.Organizations {
		for _, orderer := range organization.Orderers {
			config := compose.ResolveOrdererDockerComposeFile(cm.config.Output, organization.Domain, orderer.Subdomain)

			if err := compose.RunContainerFromTheDockerComposeFile(cm.network, config); err != nil {
				return fmt.Errorf("Error when executing the orderer container for the organization %s, orderer %s: %v\n", organization.Name, orderer.Name, err)
			}
		}
	}

	return nil
}

func (cm *ContainerManager) RunPeerContainers() error {
	slog.Info("Executing peer containers")
	for _, organization := range cm.config.Organizations {
		for _, peer := range organization.Peers {
			peerFile := compose.ResolvePeerDockerComposeFile(cm.config.Output, organization.Domain, peer.Subdomain)
			if err := compose.RunContainerFromTheDockerComposeFile(cm.network, peerFile); err != nil {
				return fmt.Errorf("Error when executing the container for the organization %s, peer %s: %v\n", organization.Name, peer.Name, err)
			}
		}
	}

	return nil
}

func (cm *ContainerManager) RunCAContainers() error {
	slog.Info("Executing certificate authority containers")
	for _, organization := range cm.config.Organizations {
		config := compose.ResolveCertificateAuthorityDockerComposeFile(cm.config.Output, organization.Domain)

		if err := compose.RunContainerFromTheDockerComposeFile(cm.network, config); err != nil {
			return fmt.Errorf("Error when executing the certificate authority container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}

func (cm *ContainerManager) RunToolsContainers() error {
	slog.Info("Executing tools containers")
	for _, organization := range cm.config.Organizations {
		config := compose.ResolveToolsDockerComposeFile(cm.config.Output, organization.Domain)

		if err := compose.RunContainerFromTheDockerComposeFile(cm.network, config); err != nil {
			return fmt.Errorf("Error when executing the tool container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}

func (cm *ContainerManager) Start() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Start Certificate Authorities", cm.RunCAContainers},
		{"Start Orderers", cm.RunOrdererContainers},
		{"Start Peers", cm.RunPeerContainers},
		{"Start Tools", cm.RunToolsContainers},
	}

	for _, step := range steps {
		slog.Info("Executing step", "step", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}

func (cm *ContainerManager) StopCertificateAuthorities() error {
	slog.Info("Stopping the certificate authority containers")
	for _, organization := range cm.config.Organizations {
		config := compose.ResolveCertificateAuthorityDockerComposeFile(cm.config.Output, organization.Domain)

		if err := compose.StopContainerFromTheDockerComposeFile(cm.network, config); err != nil {
			return fmt.Errorf("Error when stopping the certificate authority container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}
