package network

import (
	"fmt"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
)

func (f *Network) RunOrdererContainers() error {
	fmt.Print("\n=========== Executing orderer containeres ===========\n")
	for _, organization := range f.config.Organizations {
		for _, orderer := range organization.Orderers {
			config := compose.ResolveOrdererDockerComposeFile(f.config.Output, organization.Domain, orderer.Subdomain)

			if err := compose.RunContainerFromTheDockerComposeFile(f.network, config); err != nil {
				return fmt.Errorf("Error when executing the orderer container for the organization %s, orderer %s: %v\n", organization.Name, orderer.Name, err)
			}
		}
	}

	return nil
}

func (f *Network) RunPeerContainers() error {
	fmt.Print("\n=========== Executing peer containeres ===========\n")
	for _, organization := range f.config.Organizations {
		for _, peer := range organization.Peers {
			peerFile := compose.ResolvePeerDockerComposeFile(f.config.Output, organization.Domain, peer.Subdomain)
			couchDBFile := compose.ResolvePeerCouchDBDockerComposeFile(f.config.Output, organization.Domain, peer.Subdomain)

			for _, config := range []string{peerFile, couchDBFile} {
				if err := compose.RunContainerFromTheDockerComposeFile(f.network, config); err != nil {
					return fmt.Errorf("Error when executing the container for the organization %s, peer %s: %v\n", organization.Name, peer.Name, err)
				}
			}
		}
	}

	return nil
}

func (f *Network) RunCAContainers() error {
	fmt.Print("\n=========== Executing certificate authority containers ===========\n")
	for _, organization := range f.config.Organizations {
		config := compose.ResolveCertificateAuthorityDockerComposeFile(f.config.Output, organization.Domain)

		if err := compose.RunContainerFromTheDockerComposeFile(f.network, config); err != nil {
			return fmt.Errorf("Error when executing the certificate authority container for the organization %s: %v\n", organization.Name, err)
		}
	}

	return nil
}
