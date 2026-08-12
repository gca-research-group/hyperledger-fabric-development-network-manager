package validate

import (
	"errors"

	"github.com/gca-research-group/fabric-network-orchestrator/pkg/spec"
)

func Config(configuration spec.Config) error {
	var validationErrors []error
	collect := func(err error) {
		if err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	collect(InvalidOutputDirectoryNameFn(configuration))
	collect(NoOrganizationFn(configuration))

	for i, organization := range configuration.Organizations {
		collect(EmptyOrganizationNameFn(organization, i))
		collect(EmptyOrganizationDomainFn(organization, i))
		collect(InvalidCertificateAuthorityPortFn(organization))

		for j, peer := range organization.Peers {
			collect(EmptyPeerNameFn(peer, j, organization.Name))
			collect(EmptyPeerSubdomainFn(peer, organization.Name))
			collect(InvalidPeerPortFn(peer, organization.Name))
		}

		for j, orderer := range organization.Orderers {
			collect(EmptyOrdererNameFn(orderer, j, organization.Name))
			collect(EmptyOrdererSubdomainFn(orderer, organization.Name))
			collect(InvalidOrdererPortFn(orderer, organization.Name))
		}
	}

	for _, profile := range configuration.Profiles {
		collect(EmptyProfileOrganizationsFn(profile))
	}

	for _, channel := range configuration.Channels {
		collect(EmptyChannelNameFn(channel))
		collect(EmptyChannelProfileFn(channel))
		collect(InvalidChannelNameFn(channel))

		for i, chaincode := range channel.Chaincodes {
			collect(EmptyChaincodeNameFn(chaincode, i))
			collect(EmptyChaincodePathFn(chaincode, i))
			collect(EmptyChaincodeVersionFn(chaincode, i))
		}
	}

	collect(UnsupportedChannelCapabilityFn(configuration))
	collect(UnsupportedApplicationCapabilityFn(configuration))
	collect(UnsupportedOrdererCapabilityFn(configuration))
	collect(ValidateTopologyFn(configuration))

	return errors.Join(validationErrors...)
}
