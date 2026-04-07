package network

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
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
	fmt.Print("\n=========== Executing orderer containers ===========\n")
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
	fmt.Print("\n=========== Executing peer containers ===========\n")
	for _, organization := range cm.config.Organizations {
		for _, peer := range organization.Peers {
			couchDBFile := compose.ResolvePeerCouchDBDockerComposeFile(cm.config.Output, organization.Domain, peer.Subdomain)
			peerFile := compose.ResolvePeerDockerComposeFile(cm.config.Output, organization.Domain, peer.Subdomain)

			for _, config := range []string{couchDBFile, peerFile} {
				if err := compose.RunContainerFromTheDockerComposeFile(cm.network, config); err != nil {
					return fmt.Errorf("Error when executing the container for the organization %s, peer %s: %v\n", organization.Name, peer.Name, err)
				}
			}
		}
	}

	return nil
}

func (cm *ContainerManager) RunCAContainers() error {
	fmt.Print("\n=========== Executing certificate authority containers ===========\n")
	for _, organization := range cm.config.Organizations {
		config := compose.ResolveCertificateAuthorityDockerComposeFile(cm.config.Output, organization.Domain)

		if err := compose.RunContainerFromTheDockerComposeFile(cm.network, config); err != nil {
			return fmt.Errorf("Error when executing the certificate authority container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}

func (cm *ContainerManager) RunToolsContainers() error {
	fmt.Print("\n=========== Executing tools containers ===========\n")
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
		fmt.Printf(">>> Step: %s\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed at step %s: %w", step.name, err)
		}
	}

	return nil
}

func (cm *ContainerManager) StopCertificateAuthorities() error {
	fmt.Print("\n=========== Stopping the certificate authority containers ===========\n")
	for _, organization := range cm.config.Organizations {
		config := compose.ResolveCertificateAuthorityDockerComposeFile(cm.config.Output, organization.Domain)

		if err := compose.StopContainerFromTheDockerComposeFile(cm.network, config); err != nil {
			return fmt.Errorf("Error when stopping the certificate authority container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}
